package service

import (
	"container/list"
	"context"
	"sort"
	"sync"
	"time"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type Server struct {
	eventctrl.UnimplementedEventsServer
	sync.RWMutex

	sessions    map[int64]map[int64]*list.Element
	eventsList  *list.List
	brokerChan  chan *models.Event
	listUpdated chan bool
}

func RunEventsService() *Server {
	pubChan := make(chan *models.Event, 1)
	publish(pubChan)

	srv := Server{
		sessions:    make(map[int64]map[int64]*list.Element),
		eventsList:  list.New(),
		brokerChan:  pubChan,
		listUpdated: make(chan bool, 1),
	}
	go srv.bypassTimer()

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
			s.passedEvents(t1)
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

func (s *Server) passedEvents(currentTime time.Time) {
	for e := s.eventsList.Front(); e != nil; e = e.Next() {
		event := e.Value.(*models.Event)
		eventTime := time.UnixMilli(event.Time).UTC()
		timeDuration := eventTime.Sub(currentTime)

		if timeDuration <= 0 {
			s.brokerChan <- event
			s.Lock()
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(e)
			s.Unlock()
			continue
		}
		return
	}
}

func (s *Server) MakeEvent(
	ctx context.Context,
	req *eventctrl.MakeEventRequest,
) (*eventctrl.EventIdAvail, error) {
	s.Lock()
	defer s.Unlock()

	event := &models.Event{
		SenderID: req.SenderId,
		ID:       int64(len(s.sessions[req.SenderId]) + 1),
		Time:     req.Time,
		Name:     req.Name,
	}

	if _, isExist := s.sessions[req.SenderId]; !isExist {
		s.sessions[req.SenderId] = make(map[int64]*list.Element)
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

	return &eventctrl.EventIdAvail{
		EventId: event.ID,
	}, nil
}

func (s *Server) GetEvent(
	ctx context.Context,
	req *eventctrl.GetEventRequest,
) (*eventctrl.Event, error) {
	s.RLock()
	defer s.RUnlock()

	if eventsByClient, isCreated := s.sessions[req.SenderId]; isCreated {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			event := eventPtr.Value.(*models.Event)

			return &eventctrl.Event{
				SenderId: event.SenderID,
				EventId:  event.ID,
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

	if eventsByClient, isExist := s.sessions[req.SenderId]; isExist {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.sessions[req.SenderId], req.EventId)

			frontItem := s.eventsList.Front().Value
			s.eventsList.Remove(eventPtr)

			if eventPtr.Value == frontItem {
				s.listUpdated <- true
			}
			return &eventctrl.EventIdAvail{
				EventId: req.EventId,
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
	if s.sessions == nil {
		return ErrEventNotFound
	}

	if eventsByClient, isExist := s.sessions[req.SenderId]; isExist {
		var eventStream []*models.Event
		for _, eventPtr := range eventsByClient {
			event := eventPtr.Value.(*models.Event)
			if req.FromTime <= event.Time && event.Time <= req.ToTime {
				eventStream = append(eventStream, event)
				foundEvent = true
			}
		}

		sort.Slice(eventStream, func(i, j int) bool {
			return eventStream[i].ID < eventStream[j].ID
		})

		for _, event := range eventStream {
			if err := stream.Send(AccumulateEvent(*event)); err != nil {
				return ErrEventNotFound
			}
		}
	}

	if !foundEvent {
		return nil
	}
	return nil
}

func AccumulateEvent(event models.Event) *eventctrl.Event {
	return &eventctrl.Event{
		SenderId: event.SenderID,
		EventId:  event.ID,
		Time:     event.Time,
		Name:     event.Name,
	}
}
