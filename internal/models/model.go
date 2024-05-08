package models

import (
	"strconv"
	"time"

	eventcrtl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
	"github.com/emersion/go-ical"
	"github.com/google/uuid"
)

type Event struct {
	ID       uuid.UUID
	SenderID int64
	Time     int64
	Name     string
}

func (e *Event) ICalObjectBuilder(req *eventcrtl.MakeEventRequest) *ical.Calendar {
	alarm := ical.NewComponent(ical.CompAlarm)
	alarm.Props.SetText(ical.PropAction, "DISPLAY")
	trigger := ical.NewProp(ical.PropTrigger)
	trigger.SetDuration(-58 * time.Minute)
	alarm.Props.Set(trigger)

	calEvent := ical.NewEvent()
	calEvent.Name = ical.CompEvent
	for k, v := range req.Details.GetFields() {
		prop := ical.NewProp(k)
		if prop.ValueType() == ical.ValueDateTime {
			calEvent.Props.SetDateTime(k, time.UnixMilli(int64(v.GetNumberValue())))
		}
		if prop.ValueType() == ical.ValueInt || k == "X-PROTEI-SENDERID" {
			calEvent.Props.SetText(k, strconv.Itoa(int(v.GetNumberValue())))
		}
		if prop.ValueType() == ical.ValueText {
			calEvent.Props.SetText(k, v.GetStringValue())
		}
	}
	calEvent.Props.SetDateTime(ical.PropDateTimeStamp, time.Now().UTC())
	calEvent.Props.SetText(ical.PropSequence, "1")
	calEvent.Props.SetText(ical.PropUID, e.ID.String())
	calEvent.Props.SetText(ical.PropSummary, e.Name)
	calEvent.Props.SetText(ical.PropTransparency, "OPAQUE")
	calEvent.Props.SetText(ical.PropClass, "PUBLIC")
	calEvent.Children = []*ical.Component{alarm}

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//Raimguzhinov//go-caldav 1.0//EN")
	cal.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
	cal.Children = []*ical.Component{calEvent.Component}

	return cal
}
