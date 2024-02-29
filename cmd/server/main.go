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
	var host string
	var port uint
	flag.StringVar(&host, "h", "localhost", "host address")
	flag.UintVar(&port, "p", 8080, "port number")
	flag.Parse()

	lis, err := net.Listen("tcp", host+":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	eventsrv := service.RunEventsService()
	eventmanager.RegisterEventsServer(s, eventsrv)
	log.Printf("Listening on %s:%d", host, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
