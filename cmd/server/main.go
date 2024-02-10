package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	eventsrv "github.com/Raimguzhinov/simple-grpc/internal/server"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
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

	server := eventsrv.RunEventsServer()
	eventmanager.RegisterEventsServer(s, server)
	log.Printf("Listening on %s:%d", *host, *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
