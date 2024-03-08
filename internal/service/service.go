package service

import (
	"container/list"
	"context"
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
	return s.sessions != nil && s.brokerChan != nil
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
		timer := time.NewTimer(timeDuration)

		select {
		case <-timer.C:
			s.brokerChan <- event
			delete(s.sessions[event.SenderID], event.ID)
			s.eventsList.Remove(eventPtr)
		case <-s.listUpdated:
			timer.Stop()
		}
	}
}

func (s *Server) MakeEvent(
	ctx context.Context,
	req *eventctrl.MakeEventRequest,
) (*eventctrl.MakeEventResponse, error) {
	s.RLock()
	event := &models.Event{
		SenderID: req.SenderId,
		ID:       int64(len(s.sessions[req.SenderId]) + 1),
		Time:     req.Time,
		Name:     req.Name,
	}
	s.RUnlock()

	if _, isCreated := s.sessions[req.SenderId]; !isCreated {
		s.Lock()
		s.sessions[req.SenderId] = make(map[int64]*list.Element)
		s.Unlock()
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

	s.Lock()
	s.sessions[req.SenderId][event.ID] = eventPtr
	s.Unlock()

	return &eventctrl.MakeEventResponse{
		EventId: event.ID,
	}, nil
}

func (s *Server) GetEvent(
	ctx context.Context,
	req *eventctrl.GetEventRequest,
) (*eventctrl.GetEventResponse, error) {
	// FIX: Add mutex lock

	if eventsByClient, isCreated := s.sessions[req.SenderId]; isCreated {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			event := eventPtr.Value.(*models.Event)
			return &eventctrl.GetEventResponse{
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
) (*eventctrl.DeleteEventResponse, error) {
	if eventsByClient, isCreated := s.sessions[req.SenderId]; isCreated {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.sessions[req.SenderId], req.EventId)

			frontItem := s.eventsList.Front().Value
			s.eventsList.Remove(eventPtr)

			if eventPtr.Value == frontItem {
				s.listUpdated <- true
			}
			return &eventctrl.DeleteEventResponse{
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
	foundEvent := false
	if s.sessions == nil {
		return ErrEventNotFound
	}
	if eventsByClient, ok := s.sessions[req.SenderId]; ok {
		for _, eventPtr := range eventsByClient {
			event := eventPtr.Value.(*models.Event)
			if req.FromTime <= event.Time && event.Time <= req.ToTime {
				foundEvent = true
				if err := stream.Send(AccumulateEvent(*event)); err != nil {
					return ErrEventNotFound
				}
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
