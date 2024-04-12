package user

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

type ClientManager interface {
	EventMaker(senderID int64, dateCer string, timeCer string, eventName string)
	EventGetter(senderID int64, eventID string)
	EventDeleter(senderID int64, eventID string)
	EventsGetter(senderID int64, dateFrom string, timeFrom string, dateTo string, timeTo string)
}

type User struct {
	eventctrl.EventsClient
}

func NewUser(conn *eventctrl.EventsClient) *User {
	return &User{
		EventsClient: *conn,
	}
}

type CmdLiner interface {
	promtPicker() string
	promtMaker() string
	promtIDGetter() string
	promtTimeRangeGetter() string
	promtExit()
}

func RunEventsClient(client eventctrl.EventsClient, senderID int64) {
	u := NewUser(&client)

	routingKey := strconv.Itoa(int(senderID))
	queueName := strconv.Itoa(int(senderID))

	f := bufio.NewWriter(os.Stdout)
	go notifyHandler(&eventctrl.Event{}, routingKey, queueName, f)

	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		switch u.promtPicker(f) {
		case "MakeEvent":
			u.EventMaker(senderID, u.promtMaker(), os.Stdout)
		case "GetEvent":
			u.EventGetter(senderID, u.promtIDGetter(), os.Stdout)
		case "DeleteEvent":
			u.EventDeleter(senderID, u.promtIDGetter(), os.Stdout)
		case "GetEvents":
			u.EventsGetter(senderID, u.promtTimeRangeGetter(), os.Stdout)
		case "exit":
			u.promtExit()
		default:
			fmt.Fprintln(os.Stderr, ErrBadProcedure)
		}
	}
}

func (u *User) EventMaker(senderID int64, promt string, w io.Writer) {
	var eventName, dateCer, timeCer string
	if _, err := fmt.Sscan(promt, &dateCer, &timeCer, &eventName); err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	locDateTime, err := time.ParseInLocation(time.DateTime, dateCer+" "+timeCer, time.Local)
	dateTime := locDateTime.UTC()
	if err != nil {
		fmt.Fprintln(w, ErrBadDateTime)
		return
	}
	m := map[string]interface{}{
		"STATUS":            "CONFIRMED",
		"X-PROTEI-SENDERID": senderID,
		"DESCRIPTION":       eventName,
		"CREATED":           time.Now().UnixMilli(),
		"LAST-MODIFIED":     time.Now().UnixMilli(),
		"DTSTART":           dateTime.UnixMilli(),
		"DTEND":             dateTime.Add(time.Hour * 24).UnixMilli(),
	}
	details, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	res, err := u.MakeEvent(context.Background(), &eventctrl.MakeEventRequest{
		SenderId: senderID,
		Time:     dateTime.UnixMilli(),
		Name:     eventName,
		Details:  details,
	})
	if err != nil {
		fmt.Fprintln(w, ErrUnexpected)
	}
	eventID, err := uuid.FromBytes(res.EventId)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "Created {%s}\n", eventID)
}

func (u *User) EventGetter(senderID int64, promt string, w io.Writer) {
	var eventStr string
	if _, err := fmt.Sscan(promt, &eventStr); err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	eventID, err := uuid.Parse(eventStr)
	if err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	binUID, _ := eventID.MarshalBinary()
	res, err := u.GetEvent(context.Background(), &eventctrl.GetEventRequest{
		SenderId: senderID,
		EventId:  binUID,
	})
	if err != nil {
		fmt.Fprintln(w, ErrNotFound)
	} else {
		eventID, err := uuid.FromBytes(res.EventId)
		if err != nil {
			return
		}
		t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
		fmt.Fprintf(w, "Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n", res.SenderId, eventID, t, res.Name)
	}
}

func (u *User) EventDeleter(senderID int64, promt string, w io.Writer) {
	var eventStr string
	if _, err := fmt.Sscan(promt, &eventStr); err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	eventID, err := uuid.Parse(eventStr)
	if err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	binUID, _ := eventID.MarshalBinary()
	res, err := u.DeleteEvent(context.Background(), &eventctrl.DeleteEventRequest{
		SenderId: senderID,
		EventId:  binUID,
	})
	if err != nil {
		parts := strings.Split(err.Error(), "desc = ")
		fmt.Fprintln(w, parts[1])
	} else {
		eventID, err := uuid.FromBytes(res.EventId)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "Deleted {%s}\n", eventID)
	}
}

func (u *User) EventsGetter(senderID int64, promt string, w io.Writer) {
	var dateFrom, dateTo, timeFrom, timeTo string
	if _, err := fmt.Sscan(promt, &dateFrom, &timeFrom, &dateTo, &timeTo); err != nil {
		fmt.Fprintln(w, ErrBadFormat)
		return
	}
	locDateTimeFrom, err1 := time.ParseInLocation(time.DateTime, dateFrom+" "+timeFrom, time.Local)
	locDateTimeTo, err2 := time.ParseInLocation(time.DateTime, dateTo+" "+timeTo, time.Local)
	dateTimeFrom := locDateTimeFrom.UTC()
	dateTimeTo := locDateTimeTo.UTC()

	if err1 != nil || err2 != nil {
		fmt.Fprintln(w, ErrBadDateTime)
		return
	}
	stream, err := u.GetEvents(context.Background(), &eventctrl.GetEventsRequest{
		SenderId: senderID,
		FromTime: dateTimeFrom.UnixMilli(),
		ToTime:   dateTimeTo.UnixMilli(),
	})
	if err != nil {
		fmt.Fprintln(w, ErrUnexpected)
		return
	}
	for i := 0; ; i++ {
		res, err := stream.Recv()
		if err == io.EOF {
			if i == 0 {
				fmt.Fprintln(w, ErrNotFound)
			}
			break
		}
		eventID, err := uuid.FromBytes(res.EventId)
		if err != nil {
			fmt.Fprintln(w, ErrNotFound)
			return
		}
		t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
		fmt.Fprintf(w, "Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n", res.SenderId, eventID, t, res.Name)
	}
}

func (u *User) promtPicker(f *bufio.Writer) string {
	f.Flush()
	sp := selection.New(
		"\nWhat do you pick?",
		[]string{"MakeEvent", "GetEvent", "DeleteEvent", "GetEvents", "exit"},
	)
	sp.PageSize = 4
	choice, err := sp.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	return choice
}

func (u *User) promtExit() {
	input := confirmation.New("Are you confirming the exit?", confirmation.Yes)
	ready, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if ready {
		os.Exit(0)
	}
}

func (u *User) promtMaker() string {
	input := textinput.New("Enter <date> <time> <event_name>:")
	input.Placeholder = "Args cannot be empty"
	input.Validate = func(input string) error {
		if len(input) < 21 {
			return fmt.Errorf("few arguments")
		}
		return nil
	}
	args, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	return args
}

func (u *User) promtIDGetter() string {
	input := textinput.New("Enter <event_id>:")
	input.Placeholder = "Arg cannot be empty"
	arg, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	return arg
}

func (u *User) promtTimeRangeGetter() string {
	input := textinput.New("Enter <from_date> <from_time> <to_date> <to_time>:")
	input.Placeholder = "Args cannot be empty"
	input.Validate = func(input string) error {
		if len(input) <= 38 || len(input) > 39 {
			return fmt.Errorf("few arguments")
		}
		return nil
	}
	args, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	return args
}
