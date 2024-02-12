package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type events struct {
	ID       int64
	SenderID int64
	Time     int64
	Name     string
}

type server struct {
	eventmanager.UnimplementedEventsServer
	eventsByClient map[int64]map[int64]events
}

func RunEventsService() *server {
	return &server{
		eventsByClient: make(map[int64]map[int64]events),
	}
}

func (s *server) MakeEvent(ctx context.Context, req *eventmanager.MakeEventRequest) (*eventmanager.MakeEventResponse, error) {
	event := events{
		SenderID: req.SenderId,
		Time:     req.Time,
		Name:     req.Name,
	}
	if _, isCreated := s.eventsByClient[req.SenderId]; !isCreated {
		s.eventsByClient[req.SenderId] = make(map[int64]events)
	}
	event.ID = int64(len(s.eventsByClient[req.SenderId]) + 1)
	s.eventsByClient[req.SenderId][event.ID] = event
	publish(event)

	return &eventmanager.MakeEventResponse{
		EventId: event.ID,
	}, nil
}

func (s *server) GetEvent(ctx context.Context, req *eventmanager.GetEventRequest) (*eventmanager.GetEventResponse, error) {
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
	return nil, errors.New("Not found")
}

func (s *server) DeleteEvent(ctx context.Context, req *eventmanager.DeleteEventRequest) (*eventmanager.DeleteEventResponse, error) {
	if eventsByClient, isCreated := s.eventsByClient[req.SenderId]; isCreated {
		if _, isCreated := eventsByClient[req.EventId]; isCreated {
			delete(s.eventsByClient[req.SenderId], req.EventId)
			return &eventmanager.DeleteEventResponse{
				EventId: req.EventId,
			}, nil
		}
	}
	return nil, errors.New("Not found")
}

func (s *server) GetEvents(req *eventmanager.GetEventsRequest, stream eventmanager.Events_GetEventsServer) error {
	senderID := req.SenderId
	for _, eventsByClient := range s.eventsByClient {
		for _, event := range eventsByClient {
			if event.SenderID == senderID {
				if req.FromTime < event.Time && event.Time < req.ToTime {
					if err := stream.Send(accumulateEvent(event)); err != nil {
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

func accumulateEvent(event events) *eventmanager.Event {
	return &eventmanager.Event{
		SenderId: event.SenderID,
		EventId:  event.ID,
		Time:     event.Time,
		Name:     event.Name,
	}
}

func publish(event events) {
	exchange := "event.ex"
	routingKey := strconv.Itoa(int(event.SenderID))
	queueName := strconv.Itoa(int(event.SenderID))

	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			panic(err)
		}
		defer ch.Close()

		if err := ch.ExchangeDeclare(
			exchange, // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		); err != nil {
			panic(err)
		}

		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // noWait
			nil,       // arguments
		)
		if err != nil {
			panic(err)
		}

		if err := ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange
			false,      // noWait
			nil,        // arguments
		); err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		t1 := time.Now().Local()
		t2 := time.UnixMilli(event.Time)
		timeDuration := t2.Sub(t1)
		if timeDuration < 0 {
			return
		}
		timer := time.NewTimer(timeDuration)
		<-timer.C
		body, err := proto.Marshal(accumulateEvent(event))
		if err != nil {
			panic(err)
		}

		err = ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
			Body: body,
		})
		if err != nil {
			panic(err)
		}

		log.Println("[x] Sent:", body)
	}()
}