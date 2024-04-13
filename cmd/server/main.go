package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/Raimguzhinov/simple-grpc/configs"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	var host string
	var port uint
	flag.StringVar(&host, "host", cfg.EventServer.Host, "host address")
	flag.UintVar(&port, "p", cfg.EventServer.Port, "port number")
	flag.Parse()

	lis, err := net.Listen("tcp", host+":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Fatalf("Server failed to start due to err: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	cal := service.NewCalendarService(
		cfg.CaldavServer.Url,
		cfg.CaldavServer.Login,
		cfg.CaldavServer.Password,
	)
	eventsrv := service.RunEventsService()
	eventsrv.RegisterCalDAVServer(cal)
	eventctrl.RegisterEventsServer(s, eventsrv)
	log.Printf("Listening on %s:%d", host, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
