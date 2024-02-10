package main

import (
	"flag"
	"log"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	eventclt "github.com/Raimguzhinov/simple-grpc/internal/client"
	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

func main() {
	remote := flag.String("dst", "localhost", "remote address")
	port := flag.Int("p", 8080, "port number")
	senderID := flag.Int64("sender-id", 1, "sender id")
	flag.Parse()

	var wg sync.WaitGroup

	conn, err := grpc.Dial(*remote+":"+strconv.Itoa(*port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := eventmanager.NewEventsClient(conn)

	wg.Add(1)
	go func() {
		defer wg.Done()
		eventclt.RunEventsClient(client, senderID)
	}()
	wg.Wait()
}
