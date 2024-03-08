package main

import (
	"flag"
	"log"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Raimguzhinov/simple-grpc/internal/user"
	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func main() {
	var remote string
	var port uint
	var senderID int64
	flag.StringVar(&remote, "dst", "localhost", "remote address")
	flag.UintVar(&port, "p", 8080, "port number")
	flag.Int64Var(&senderID, "sender-id", 1, "sender id")
	flag.Parse()

	var wg sync.WaitGroup

	conn, err := grpc.Dial(
		remote+":"+strconv.Itoa(int(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := eventctrl.NewEventsClient(conn)

	wg.Add(1)
	go func() {
		defer wg.Done()
		user.RunEventsClient(client, senderID)
	}()
	wg.Wait()
}
