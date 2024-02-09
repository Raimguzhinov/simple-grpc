package eventsrv

import (
	"context"
	"errors"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

type events struct {
	ID       int64
	SenderID int64
	Time     int64
	Name     string
}

type Server struct {
	eventmanager.UnimplementedEventsServer
	EventsByClient map[int64]map[int64]events
}

func (s *Server) MakeEvent(ctx context.Context, req *eventmanager.MakeEventRequest) (*eventmanager.MakeEventResponse, error) {
	event := events{
		SenderID: req.SenderId,
		Time:     req.Time,
		Name:     req.Name,
	}
	if _, isCreated := s.EventsByClient[req.SenderId]; !isCreated {
		s.EventsByClient[req.SenderId] = make(map[int64]events)
	}
	event.ID = int64(len(s.EventsByClient[req.SenderId]) + 1)
	s.EventsByClient[req.SenderId][event.ID] = event
	return &eventmanager.MakeEventResponse{
		EventId: event.ID,
	}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *eventmanager.GetEventRequest) (*eventmanager.GetEventResponse, error) {
	if eventsByClient, isCreated := s.EventsByClient[req.SenderId]; isCreated {
		if event, isCreated := eventsByClient[req.EventId]; isCreated {
			return &eventmanager.GetEventResponse{
				SenderId: event.SenderID,
				EventId:  event.ID,
				Time:     event.Time,
				Name:     event.Name,
			}, nil
		}
	}
	return nil, errors.New("Not found")
}

func (s *Server) DeleteEvent(ctx context.Context, req *eventmanager.DeleteEventRequest) (*eventmanager.DeleteEventResponse, error) {
	if eventsByClient, isCreated := s.EventsByClient[req.SenderId]; isCreated {
		if _, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.EventsByClient[req.SenderId], req.EventId)
			return &eventmanager.DeleteEventResponse{
				EventId: req.EventId,
			}, nil
		}
	}
	return nil, errors.New("Not found")
}

func (s *Server) GetEvents(req *eventmanager.GetEventsRequest, stream eventmanager.Events_GetEventsServer) error {
	senderID := req.SenderId
	for _, eventsByClient := range s.EventsByClient {
		for _, event := range eventsByClient {
			if event.SenderID == senderID {
				if req.FromTime < event.Time && event.Time < req.ToTime {
					if err := stream.Send(&eventmanager.GetEventsResponse{
						SenderId: event.SenderID,
						EventId:  event.ID,
						Time:     event.Time,
						Name:     event.Name,
					}); err != nil {
						return err
					}
				} else {
					return errors.New("Not found")
				}
			} else {
				return errors.New("Not found")
			}
		}
	}
	return nil
}

func NewEventsServer() *Server {
	return &Server{
		EventsByClient: make(map[int64]map[int64]events),
	}
}
