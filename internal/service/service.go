package service

import (
	"container/list"
	"context"
	"log"
	"sync"
	"time"

	"github.com/Raimguzhinov/go-webdav/caldav"
	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
	"github.com/google/uuid"
)

type Server struct {
	eventctrl.UnimplementedEventsServer
	sync.RWMutex

	calDAVServer *Calendar
	calendars    []caldav.Calendar
	sessions     map[int64]map[uuid.UUID]*list.Element
	eventsList   *list.List
	brokerChan   chan *models.Event
	listUpdated  chan bool
}

func RunEventsService(events ...*eventctrl.MakeEventRequest) *Server {
	pubChan := make(chan *models.Event, 10000)
	publish(pubChan)

	srv := Server{
		sessions:    make(map[int64]map[uuid.UUID]*list.Element),
		eventsList:  list.New(),
		brokerChan:  pubChan,
		listUpdated: make(chan bool, 1),
	}
	for _, event := range events {
		_, _ = srv.MakeEvent(context.Background(), event)
	}
	go srv.bypassTimer()

	return &srv
}

func (s *Server) RegisterCalDAVServer(calDAVServer *Calendar) {
	calendars, err := calDAVServer.GetCalendars(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	s.calDAVServer = calDAVServer
	s.calendars = calendars
	go s.syncWithCalendars()
}

func (s *Server) IsInitialized() bool {
	return s.eventsList != nil && s.sessions != nil && s.brokerChan != nil
}

func (s *Server) syncWithCalendars() {
	s.Lock()
	defer s.Unlock()

	var calEvents []*models.Event
	var err error
	for _, calendar := range s.calendars {
		calEvents, err = s.calDAVServer.LoadEvents(context.Background(), calendar)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, req := range calEvents {
		if _, isExist := s.sessions[req.SenderID]; !isExist {
			s.sessions[req.SenderID] = make(map[uuid.UUID]*list.Element)
		}

		var eventPtr *list.Element
		if s.eventsList.Len() == 0 {
			eventPtr = s.eventsList.PushBack(req)
			s.listUpdated <- true
		} else {
			for e := s.eventsList.Back(); e != nil; e = e.Prev() {
				item := e.Value.(*models.Event)
				if req.Time >= item.Time {
					eventPtr = s.eventsList.InsertAfter(req, e)
					break
				} else if e == s.eventsList.Front() && item.Time > req.Time {
					eventPtr = s.eventsList.InsertBefore(req, e)
					s.listUpdated <- true
					break
				}
			}
		}
		s.sessions[req.SenderID][req.ID] = eventPtr
	}
}

func (s *Server) bypassTimer() {
	for {
		s.Lock()
		if s.eventsList.Len() == 0 {
			s.Unlock()
			<-s.listUpdated
			continue
		}
		eventPtr := s.eventsList.Front()
		event := eventPtr.Value.(*models.Event)
		t1 := time.Now().UTC()
		t2 := time.UnixMilli(event.Time).UTC()
		timeDuration := t2.Sub(t1)
		if timeDuration <= 0 {
			s.Unlock()
			_ = s.PassedEvents(t1)
			continue
		}
		timer := time.NewTimer(timeDuration)
		s.Unlock()

		select {
		case <-timer.C:
			s.Lock()
			s.brokerChan <- event
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(eventPtr)
			s.Unlock()
		case <-s.listUpdated:
			timer.Stop()
		}
	}
}

func (s *Server) PassedEvents(currentTime time.Time) int {
	s.Lock()
	defer s.Unlock()
	for e := s.eventsList.Front(); e != nil; {
		event := e.Value.(*models.Event)
		eventTime := time.UnixMilli(event.Time).UTC()
		timeDuration := eventTime.Sub(currentTime)

		if timeDuration <= 0 {
			// TODO: Push old events to broker again?
			// s.brokerChan <- event
			next := e.Next()
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(e)
			e = next
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

	if s.calDAVServer != nil && len(s.calendars) > 0 {
		go func() {
			err := s.calDAVServer.PutCalendarObject(
				context.Background(),
				event.ICalObjectBuilder(req),
				s.calendars[0],
			)
			if err != nil {
				log.Println("Can't put event to calendar", err)
			}
		}()
	}

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

			if s.calDAVServer != nil && len(s.calendars) > 0 {
				go func() {
					err = s.calDAVServer.DeleteCalendarObject(
						context.Background(),
						eventID,
						s.calendars[0],
					)
					if err != nil {
						log.Println("Can't delete event from calendar", err)
					}
				}()
			}

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
