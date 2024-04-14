package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Raimguzhinov/go-webdav"
	"github.com/Raimguzhinov/go-webdav/caldav"
	"github.com/Raimguzhinov/simple-grpc/internal/models"
	"github.com/emersion/go-ical"
	"github.com/google/uuid"
)

type Calendar struct {
	caldavUrl      string
	caldavLogin    string
	caldavPassword string
}

type transport struct {
	current *http.Request
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	// fmt.Printf("\n\nRequest:\n%v\n", req)
	resp, err := http.DefaultTransport.RoundTrip(req)
	return resp, err
}

func NewCalendarService(caldavUrl string, caldavLogin string, caldavPassword string) *Calendar {
	return &Calendar{
		caldavUrl:      caldavUrl,
		caldavLogin:    caldavLogin,
		caldavPassword: caldavPassword,
	}
}

func (c *Calendar) getClient() (*caldav.Client, error) {
	httpClient := webdav.HTTPClientWithBasicAuth(
		&http.Client{Transport: &transport{}},
		c.caldavLogin,
		c.caldavPassword,
	)
	client, err := caldav.NewClient(httpClient, c.caldavUrl)
	return client, err
}

func (c *Calendar) GetCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	client, err := c.getClient()
	if err != nil {
		log.Println("Error get client for principal", &c.caldavLogin, err)
		return make(
				[]caldav.Calendar,
				0,
			), errors.New(
				fmt.Sprintf("Error get client in principal method for user %s", c.caldavLogin),
			)
	}
	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		log.Println("Error get principal", &c.caldavLogin, err)
		return make(
				[]caldav.Calendar,
				0,
			), errors.New(
				fmt.Sprintf("Error get principal for user %s", c.caldavLogin),
			)
	}
	calendarHomeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		log.Println("Error get calendarHomeSet", &c.caldavLogin, err)
		return make(
				[]caldav.Calendar,
				0,
			), errors.New(
				fmt.Sprintf("Error get calendarHomeSet for user %s", c.caldavLogin),
			)
	}
	return client.FindCalendars(ctx, calendarHomeSet)
}

func (c *Calendar) LoadEvents(
	ctx context.Context,
	calendar caldav.Calendar,
) ([]*models.Event, error) {
	var events []*models.Event
	client, err := c.getClient()
	if err != nil {
		log.Println("Error get client for principal", &c.caldavLogin, err)
		return events, errors.New(
			fmt.Sprintf("Error get client in principal method for user %s", c.caldavLogin),
		)
	}
	calendarObjects, err := c.queryCalendarEvents(ctx, client, calendar.Path)
	if calendarObjects == nil {
		log.Println("Can't get events for calendar "+calendar.Path, err)
		return events, errors.New("Can't get events from calendar")
	}
	timezone, err := GetTimezone(calendarObjects)
	if err != nil {
		log.Println("Can't get timezone for calendar "+calendar.Path, err)
	}
	eventDtos, err := CalendarObjectToEventArray(calendarObjects, timezone)
	if err != nil {
		log.Println("Can't parse events for calendar "+calendar.Path, err)
		return events, errors.New("Can't parse events from calendar")
	}
	events = append(events, eventDtos...)
	return events, nil
}

func (c *Calendar) PutCalendarObject(
	ctx context.Context,
	cal *ical.Calendar,
	calendars ...caldav.Calendar,
) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}
	for _, calendar := range calendars {
		eventID := getPropertyValue(cal.Children[0].Props.Get(ical.PropUID))
		obj, err := client.PutCalendarObject(ctx, calendar.Path+eventID+".ics", cal)
		if err != nil {
			return err
		}
		log.Println("Putted:", obj)
	}
	return nil
}

func (c *Calendar) DeleteCalendarObject(
	ctx context.Context,
	eventID uuid.UUID,
	calendars ...caldav.Calendar,
) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}
	for _, calendar := range calendars {
		err = client.RemoveAll(ctx, calendar.Path+eventID.String()+".ics")
		if err != nil {
			return err
		}
		log.Println("Deleted:", calendar.Path+eventID.String()+".ics")
	}
	return nil
}

func (c *Calendar) queryCalendarEvents(
	ctx context.Context,
	client *caldav.Client,
	calendarPath string,
) ([]caldav.CalendarObject, error) {
	query := &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name: "VEVENT",
			}},
		},
	}
	return client.QueryCalendar(ctx, calendarPath, query)
}

func CalendarObjectToEventArray(
	calendarObjects []caldav.CalendarObject,
	timezone string,
) ([]*models.Event, error) {
	eventById := make(map[uuid.UUID]*models.Event)
	for _, calendarObject := range calendarObjects {
		for _, e := range calendarObject.Data.Events() {
			eventSenderId, _ := strconv.Atoi(getPropertyValue(e.Props.Get("X-PROTEI-SENDERID")))
			eventId := uuid.MustParse(getPropertyValue(e.Props.Get("UID")))
			if _, ok := eventById[eventId]; ok {
				continue
			}
			location, _ := time.LoadLocation(timezone)
			eventName, _ := e.Props.Text("SUMMARY")
			eventClass := getPropertyValue(e.Props.Get("CLASS"))
			eventDescription, _ := e.Props.Text("DESCRIPTION")
			eventUrl := getPropertyValue(e.Props.Get("URL"))
			sequence := getPropertyValue(e.Props.Get("SEQUENCE"))
			transp := getPropertyValue(e.Props.Get("TRANSP"))
			status, _ := e.Status()

			createdTime, err := e.Props.DateTime("CREATED", location)
			if err != nil {
				return nil, errors.Join(err, errors.New("Can't parse CREATED for event "+eventName))
			}

			startTime, err := e.Props.DateTime("DTSTART", location)
			if err != nil {
				return nil, errors.Join(err, errors.New("Can't parse DTSTART for event "+eventName))
			}

			endTime, err := e.Props.DateTime("DTEND", location)
			if err != nil {
				return nil, errors.Join(err, errors.New("Can't parse DTEND for event "+eventName))
			}

			dtstamp, err := e.Props.DateTime("DTSTAMP", location)
			if err != nil {
				return nil, errors.Join(err, errors.New("Can't parse DTSTAMP for event "+eventName))
			}

			lastModifiedTime, err := e.Props.DateTime("LAST-MODIFIED", location)
			if err != nil {
				return nil, errors.Join(
					err,
					errors.New("Can't parse LAST-MODIFIED for event "+eventName),
				)
			}

			log.Printf(
				"\nCLASS: %s\nSTATUS: %s\nUID: %s\nSUMMARY: %s\nDESCRIPTION: %s\nURL: %s\nCREATED: %s\nDTSTART: %s\nDTEND: %s\nDTSTAMP %s\nLAST-MODIFIED: %s\nSEQUENCE: %s\nTRANSP: %s\nX-PROTEI-SENDERID: %d\n\n",
				eventClass,
				status,
				eventId,
				eventName,
				eventDescription,
				eventUrl,
				createdTime.Local(),
				startTime,
				endTime,
				dtstamp.Local(),
				lastModifiedTime.Local(),
				sequence,
				transp,
				eventSenderId,
			)
			eventById[eventId] = &models.Event{
				ID:       eventId,
				SenderID: int64(eventSenderId),
				Time:     startTime.Unix(),
				Name:     eventName,
			}
		}
	}
	events := make([]*models.Event, 0, len(eventById))
	for _, event := range eventById {
		events = append(events, event)
	}
	return events, nil
}

func GetTimezone(calendarObjects []caldav.CalendarObject) (string, error) {
	if len(calendarObjects) == 0 {
		return "Etc/UTC", nil
	}
	for _, calendarObject := range calendarObjects {
		for _, child := range calendarObject.Data.Children {
			if child.Name == ical.CompTimezone {
				return child.Props.Text("TZID")
			}
		}
	}
	return "Etc/UTC", errors.New("Timezone not found")
}

func getPropertyValue(prop *ical.Prop) string {
	if prop == nil {
		return ""
	}
	return prop.Value
}
