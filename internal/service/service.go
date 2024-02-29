package service

import (
	"container/list"
	"context"
	"time"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type Server struct {
	eventmanager.UnimplementedEventsServer
	eventsByClient      map[int64]map[int64]*list.Element
	eventsList          *list.List
	eventsChan          chan models.Events
	listChangingChannel chan bool
}

func RunEventsService() *Server {
	pubChan := make(chan models.Events, 1)
	publish(pubChan)
	srv := Server{
		eventsByClient:      make(map[int64]map[int64]*list.Element),
		eventsList:          list.New(),
		eventsChan:          pubChan,
		listChangingChannel: make(chan bool),
	}
	go srv.timerQueue()
	return &srv
}

func (s *Server) IsInitialized() bool {
	return s.eventsByClient != nil && s.eventsChan != nil
}

func (s *Server) timerQueue() {
ResetTimer:
	for {
		if s.eventsList.Len() > 0 {
			eventPtr := s.eventsList.Front()
			event := eventPtr.Value.(models.Events)
			t1 := time.Now().UTC()
			t2 := time.UnixMilli(event.Time).UTC()
			timeDuration := t2.Sub(t1)
			timer := time.NewTimer(timeDuration)

			select {
			case <-timer.C:
				s.eventsChan <- event
				delete(s.eventsByClient[event.SenderID], event.ID)
				s.eventsList.Remove(eventPtr)
				continue ResetTimer
			case <-s.listChangingChannel:
				timer.Stop()
				continue ResetTimer
			}
		}
	}
}

func (s *Server) MakeEvent(ctx context.Context, req *eventmanager.MakeEventRequest) (*eventmanager.MakeEventResponse, error) {
	event := models.Events{
		SenderID: req.SenderId,
		Time:     req.Time,
		Name:     req.Name,
	}
	event.ID = int64(len(s.eventsByClient[req.SenderId]) + 1)
	var eventPtr *list.Element
	if _, isCreated := s.eventsByClient[req.SenderId]; !isCreated {
		s.eventsByClient[req.SenderId] = make(map[int64]*list.Element)
	}
	if s.eventsList.Len() == 0 {
		eventPtr = s.eventsList.PushBack(event)
	} else {
		for e := s.eventsList.Back(); e != nil; e = e.Prev() {
			item := models.Events(e.Value.(models.Events))
			if event.Time >= item.Time {
				eventPtr = s.eventsList.InsertAfter(event, e)
				break
			} else if e == s.eventsList.Front() && item.Time > event.Time {
				eventPtr = s.eventsList.InsertBefore(event, e)
				s.listChangingChannel <- true
				break
			}
		}
	}
	s.eventsByClient[req.SenderId][event.ID] = eventPtr
	return &eventmanager.MakeEventResponse{
		EventId: event.ID,
	}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *eventmanager.GetEventRequest) (*eventmanager.GetEventResponse, error) {
	if eventsByClient, isCreated := s.eventsByClient[req.SenderId]; isCreated {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			event := eventPtr.Value.(models.Events)
			return &eventmanager.GetEventResponse{
				SenderId: event.SenderID,
				EventId:  event.ID,
				Time:     event.Time,
				Name:     event.Name,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

func (s *Server) DeleteEvent(ctx context.Context, req *eventmanager.DeleteEventRequest) (*eventmanager.DeleteEventResponse, error) {
	if eventsByClient, isCreated := s.eventsByClient[req.SenderId]; isCreated {
		if eventPtr, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.eventsByClient[req.SenderId], req.EventId)
			frontItem := s.eventsList.Front().Value
			s.eventsList.Remove(eventPtr)
			if eventPtr.Value == frontItem {
				s.listChangingChannel <- true
			}
			return &eventmanager.DeleteEventResponse{
				EventId: req.EventId,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

func (s *Server) GetEvents(req *eventmanager.GetEventsRequest, stream eventmanager.Events_GetEventsServer) error {
	foundEvent := false
	if s.eventsByClient == nil {
		return ErrEventNotFound
	}
	if eventsByClient, ok := s.eventsByClient[req.SenderId]; ok {
		for _, eventPtr := range eventsByClient {
			event := eventPtr.Value.(models.Events)
			if req.FromTime <= event.Time && event.Time <= req.ToTime {
				foundEvent = true
				if err := stream.Send(AccumulateEvent(event)); err != nil {
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

func AccumulateEvent(event models.Events) *eventmanager.Event {
	return &eventmanager.Event{
		SenderId: event.SenderID,
		EventId:  event.ID,
		Time:     event.Time,
		Name:     event.Name,
	}
}
