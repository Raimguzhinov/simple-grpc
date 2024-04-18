package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Raimguzhinov/simple-grpc/configs"
	"github.com/Raimguzhinov/simple-grpc/internal/service"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func main() {
	cfg, err := configs.New("config.yml")
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
	broker := service.NewBroker(
		cfg.RabbitmqServer.Url,
		cfg.RabbitmqServer.Login,
		cfg.RabbitmqServer.Password,
		cfg.RabbitmqServer.Exchange,
		cfg.RabbitmqServer.ReconnDelay,
	)

	eventsrv := service.RunEventsService()
	eventsrv.RegisterCalDAVServer(cal)
	eventsrv.RegisterBrokerServer(broker)

	eventctrl.RegisterEventsServer(s, eventsrv)

	log.Printf("Listening on %s:%d", host, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
