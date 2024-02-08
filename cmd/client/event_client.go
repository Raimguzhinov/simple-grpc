package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

var (
	remote        *string
	port          *int
	senderID      *int64
	ProcedureName string
	EventID       int64
	EventName     string
	Time          int64
	TimeFrom      int64
	TimeTo        int64
)

func init() {
	remote = flag.String("dst", "localhost", "remote address")
	port = flag.Int("p", 8080, "port number")
	senderID = flag.Int64("sender-id", 1, "sender id")
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*remote+":"+strconv.Itoa(*port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		log.Fatal("Usage: -dst <remote_adress> -p <port> -sender-id <sender_id>")
	}
	defer conn.Close()
	client := eventmanager.NewEventsClient(conn)

	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		fmt.Print("> ")
		fmt.Scan(&ProcedureName)
		switch ProcedureName {
		case "MakeEvent":
			fmt.Print("Enter <time> <event_name>: ")
			fmt.Scan(&Time, &EventName)
			res, err := client.MakeEvent(context.Background(), &eventmanager.MakeEventRequest{
				SenderId: *senderID,
				Time:     Time,
				Name:     EventName,
			})
			if err != nil {
				log.Fatalf("Failed to make event: %v", err)
			}
			fmt.Println("Created {", res, "}")
		case "GetEvent":
			fmt.Print("Enter <event_id>: ")
			fmt.Scan(&EventID)
			res, err := client.GetEvent(context.Background(), &eventmanager.GetEventRequest{
				SenderId: *senderID,
				EventId:  EventID,
			})
			if err != nil {
				parts := strings.Split(err.Error(), "desc = ")
				fmt.Println(parts[1])
			} else {
				fmt.Println("Event {", res, "}")
			}
		case "DeleteEvent":
			fmt.Print("Enter <event_id>: ")
			fmt.Scan(&EventID)
			res, err := client.DeleteEvent(context.Background(), &eventmanager.DeleteEventRequest{
				SenderId: *senderID,
				EventId:  EventID,
			})
			if err != nil {
				parts := strings.Split(err.Error(), "desc = ")
				fmt.Println(parts[1])
			} else {
				fmt.Println("Deleted {", res, "}")
			}
		case "GetEvents":
			fmt.Print("Enter <from_time> <to_time>: ")
			fmt.Scan(&TimeFrom, &TimeTo)
			stream, err := client.GetEvents(context.Background(), &eventmanager.GetEventsRequest{
				SenderId: *senderID,
				FromTime: TimeFrom,
				ToTime:   TimeTo,
			})
			if err != nil {
				log.Fatalf("Failed to get events: %v", err)
			} else {
				for {
					res, err := stream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						parts := strings.Split(err.Error(), "desc = ")
						fmt.Println(parts[1])
						break
					} else {
						fmt.Println("Event {", res, "}")
					}
				}
			}
		}
	}
}
