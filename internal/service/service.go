package service

import (
	"context"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type Server struct {
	eventmanager.UnimplementedEventsServer
	eventsByClient map[int64]map[int64]models.Events
	eventsChan     chan models.Events
}

func RunEventsService() *Server {
	pubChan := make(chan models.Events, 1)
	publish(pubChan)
	return &Server{
		eventsByClient: make(map[int64]map[int64]models.Events),
		eventsChan:     pubChan,
	}
}

func (s *Server) MakeEvent(ctx context.Context, req *eventmanager.MakeEventRequest) (*eventmanager.MakeEventResponse, error) {
	event := models.Events{
		SenderID: req.SenderId,
		Time:     req.Time,
		Name:     req.Name,
	}
	if _, isCreated := s.eventsByClient[req.SenderId]; !isCreated {
		s.eventsByClient[req.SenderId] = make(map[int64]models.Events)
	}
	event.ID = int64(len(s.eventsByClient[req.SenderId]) + 1)
	s.eventsByClient[req.SenderId][event.ID] = event
	go func() {
		s.eventsChan <- event
	}()
	return &eventmanager.MakeEventResponse{
		EventId: event.ID,
	}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *eventmanager.GetEventRequest) (*eventmanager.GetEventResponse, error) {
	if eventsByClient, isCreated := s.eventsByClient[req.SenderId]; isCreated {
		if event, isCreated := eventsByClient[req.EventId]; isCreated {
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
		if _, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.eventsByClient[req.SenderId], req.EventId)
			return &eventmanager.DeleteEventResponse{
				EventId: req.EventId,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

func (s *Server) GetEvents(req *eventmanager.GetEventsRequest, stream eventmanager.Events_GetEventsServer) error {
	senderID := req.SenderId
	for _, eventsByClient := range s.eventsByClient {
		for _, event := range eventsByClient {
			if event.SenderID == senderID {
				if req.FromTime < event.Time && event.Time < req.ToTime {
					if err := stream.Send(accumulateEvent(event)); err != nil {
						return err
					}
				} else {
					return ErrEventNotFound
				}
			} else {
				return ErrEventNotFound
			}
		}
	}
	return nil
}

func accumulateEvent(event models.Events) *eventmanager.Event {
	return &eventmanager.Event{
		SenderId: event.SenderID,
		EventId:  event.ID,
		Time:     event.Time,
		Name:     event.Name,
	}
}
