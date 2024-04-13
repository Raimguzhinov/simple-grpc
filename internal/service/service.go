package service

import (
	"container/list"
	"context"
	"fmt"
	"github.com/Raimguzhinov/simple-grpc/configs"
	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
	"github.com/emersion/go-ical"
	"github.com/google/uuid"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	eventctrl.UnimplementedEventsServer
	sync.RWMutex

	calendar    *Calendar
	sessions    map[int64]map[uuid.UUID]*list.Element
	eventsList  *list.List
	brokerChan  chan *models.Event
	listUpdated chan bool
}

func RunEventsService(events ...*eventctrl.MakeEventRequest) *Server {
	cfg, err := configs.New()
	if err != nil {
		panic(err)
	}
	pubChan := make(chan *models.Event, 10000)
	publish(pubChan)

	srv := Server{
		calendar: NewCalendarService(
			cfg.CaldavServer.Url,
			cfg.CaldavServer.Login,
			cfg.CaldavServer.Password,
		),
		sessions:    make(map[int64]map[uuid.UUID]*list.Element),
		eventsList:  list.New(),
		brokerChan:  pubChan,
		listUpdated: make(chan bool, 1),
	}
	for _, event := range events {
		_, _ = srv.MakeEvent(context.Background(), event)
	}
	go srv.bypassTimer()

	calendars, err := srv.calendar.GetCalendars(context.Background())
	if err != nil {
		panic(err)
	}
	var calEvents []*models.Event
	for _, calendar := range calendars {
		calEvents, err = srv.calendar.LoadEvents(context.Background(), calendar)
		if err != nil {
			panic(err)
		}
	}
	_ = calEvents

	return &srv
}

func (s *Server) IsInitialized() bool {
	return s.eventsList != nil && s.sessions != nil && s.brokerChan != nil
}

func (s *Server) bypassTimer() {
	for {
		if s.eventsList.Len() == 0 {
			<-s.listUpdated
			continue
		}
		eventPtr := s.eventsList.Front()
		event := eventPtr.Value.(*models.Event)

		t1 := time.Now().UTC()
		t2 := time.UnixMilli(event.Time).UTC()
		timeDuration := t2.Sub(t1)
		if timeDuration <= 0 {
			_ = s.PassedEvents(t1)
			continue
		}
		timer := time.NewTimer(timeDuration)

		select {
		case <-timer.C:
			s.brokerChan <- event
			s.Lock()
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(eventPtr)
			s.Unlock()
		case <-s.listUpdated:
			timer.Stop()
		}
	}
}

func (s *Server) PassedEvents(currentTime time.Time) int {
	for e := s.eventsList.Front(); e != nil; {
		event := e.Value.(*models.Event)
		eventTime := time.UnixMilli(event.Time).UTC()
		timeDuration := eventTime.Sub(currentTime)

		if timeDuration <= 0 {
			s.brokerChan <- event
			s.Lock()
			next := e.Next()
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(e)
			e = next
			s.Unlock()
			continue
		}
		break
	}
	return s.eventsList.Len()
}

func (s *Server) MakeEvent(
	ctx context.Context,
	req *eventctrl.MakeEventRequest,
) (*eventctrl.EventIdAvail, error) {
	s.Lock()
	defer s.Unlock()

	event := &models.Event{
		SenderID: req.SenderId,
		ID:       uuid.New(),
		Time:     req.Time,
		Name:     req.Name,
	}

	go func() {
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
		calEvent.Props.SetText(ical.PropUID, event.ID.String())
		calEvent.Props.SetText(ical.PropSummary, event.Name)
		calEvent.Props.SetText(ical.PropTransparency, "OPAQUE")
		calEvent.Props.SetText(ical.PropClass, "PUBLIC")
		calEvent.Children = []*ical.Component{alarm}
		fmt.Println(calEvent.Props)

		cal := ical.NewCalendar()
		cal.Props.SetText(ical.PropVersion, "2.0")
		cal.Props.SetText(ical.PropProductID, "-//Raimguzhinov//go-caldav 1.0//EN")
		cal.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
		cal.Children = []*ical.Component{calEvent.Component}

		_ = s.calendar.PutCalendarObject(context.Background(), cal)
	}()

	if _, isExist := s.sessions[req.SenderId]; !isExist {
		s.sessions[req.SenderId] = make(map[uuid.UUID]*list.Element)
	}

	var eventPtr *list.Element
	if s.eventsList.Len() == 0 {
		eventPtr = s.eventsList.PushBack(event)
		s.listUpdated <- true
	} else {
		for e := s.eventsList.Back(); e != nil; e = e.Prev() {
			item := e.Value.(*models.Event)
			if event.Time >= item.Time {
				eventPtr = s.eventsList.InsertAfter(event, e)
				break
			} else if e == s.eventsList.Front() && item.Time > event.Time {
				eventPtr = s.eventsList.InsertBefore(event, e)
				s.listUpdated <- true
				break
			}
		}
	}
	s.sessions[req.SenderId][event.ID] = eventPtr

	binUID, err := event.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &eventctrl.EventIdAvail{
		EventId: binUID,
	}, nil
}

func (s *Server) GetEvent(
	ctx context.Context,
	req *eventctrl.GetEventRequest,
) (*eventctrl.Event, error) {
	s.RLock()
	defer s.RUnlock()

	if _, isCreated := s.sessions[req.SenderId]; isCreated {
		eventID, err := uuid.FromBytes(req.EventId)
		if err != nil {
			return nil, ErrEventNotFound
		}
		if eventPtr, isCreated := s.sessions[req.SenderId][eventID]; isCreated {
			event := eventPtr.Value.(*models.Event)

			binUID, err := event.ID.MarshalBinary()
			if err != nil {
				return nil, err
			}
			return &eventctrl.Event{
				SenderId: event.SenderID,
				EventId:  binUID,
				Time:     event.Time,
				Name:     event.Name,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

func (s *Server) DeleteEvent(
	ctx context.Context,
	req *eventctrl.DeleteEventRequest,
) (*eventctrl.EventIdAvail, error) {
	s.Lock()
	defer s.Unlock()

	if _, isExist := s.sessions[req.SenderId]; isExist {
		eventID, err := uuid.FromBytes(req.EventId)
		if err != nil {
			return nil, ErrEventNotFound
		}
		if eventPtr, isCreated := s.sessions[req.SenderId][eventID]; isCreated {
			delete(s.sessions[req.SenderId], eventID)

			frontItem := s.eventsList.Front().Value
			s.eventsList.Remove(eventPtr)

			if eventPtr.Value == frontItem {
				s.listUpdated <- true
			}

			binUID, err := eventID.MarshalBinary()
			if err != nil {
				return nil, err
			}
			return &eventctrl.EventIdAvail{
				EventId: binUID,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

func (s *Server) GetEvents(
	req *eventctrl.GetEventsRequest,
	stream eventctrl.Events_GetEventsServer,
) error {
	s.Lock()
	defer s.Unlock()
	foundEvent := false

	if _, isExist := s.sessions[req.SenderId]; isExist {
		for _, eventPtr := range s.sessions[req.SenderId] {
			event := eventPtr.Value.(*models.Event)
			if req.FromTime <= event.Time && event.Time <= req.ToTime {
				foundEvent = true
				if err := stream.Send(AccumulateEvent(*event)); err != nil {
					return ErrEventNotFound
				}
			}
		}
		if !foundEvent {
			return ErrEventNotFound
		}
	} else {
		return ErrEventNotFound
	}
	return nil
}

func (s *Server) IsEventExist(senderID int64, binUID []byte) bool {
	s.RLock()
	defer s.RUnlock()
	if eventsByClient, isCreated := s.sessions[senderID]; isCreated {
		eventID, _ := uuid.FromBytes(binUID)
		if _, isExist := eventsByClient[eventID]; isCreated {
			return isExist
		}
	}
	return false
}

func AccumulateEvent(event models.Event) *eventctrl.Event {
	binUID, _ := event.ID.MarshalBinary()
	return &eventctrl.Event{
		SenderId: event.SenderID,
		EventId:  binUID,
		Time:     event.Time,
		Name:     event.Name,
	}
}
