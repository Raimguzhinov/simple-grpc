package event_server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

type server struct {
	eventmanager.UnimplementedEventsServer
}

func (s *server) MakeEvent(ctx context.Context, req *eventmanager.MakeEventRequest) (*eventmanager.MakeEventResponse, error) {
	log.Printf("User id: %d", req.SenderId)
	return &eventmanager.MakeEventResponse{
		EventId: 1,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	eventmanager.RegisterEventsServer(s, &server{})

	log.Println("Listening on :50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
