package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func main() {
	host := flag.String("h", "localhost", "host address")
	port := flag.Int("p", 8080, "port number")
	flag.Parse()

	lis, err := net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	eventsrv := service.RunEventsService()
	eventmanager.RegisterEventsServer(s, eventsrv)
	log.Printf("Listening on %s:%d", *host, *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
