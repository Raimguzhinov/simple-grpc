package eventclt

import (
	"context"
	"fmt"
	"io"
	"strings"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

var (
	ProcedureName string
	EventID       int64
	EventName     string
	Time          int64
	TimeFrom      int64
	TimeTo        int64
)

func RunEventsClient(client eventmanager.EventsClient, senderID *int64) {
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
				return
			}
			fmt.Println("Created {", res, "}")
		case "GetEvent":
			fmt.Print("Enter <event_id>: ")
			if _, err := fmt.Scan(&EventID); err != nil {
				fmt.Println(err)
			} else {
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
			}
		case "DeleteEvent":
			fmt.Print("Enter <event_id>: ")
			if _, err := fmt.Scan(&EventID); err != nil {
				fmt.Println(err)
			} else {
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
			}
		case "GetEvents":
			fmt.Print("Enter <from_time> <to_time>: ")
			if _, err := fmt.Scan(&TimeFrom, &TimeTo); err != nil {
				fmt.Println(err)
			} else {
				stream, err := client.GetEvents(context.Background(), &eventmanager.GetEventsRequest{
					SenderId: *senderID,
					FromTime: TimeFrom,
					ToTime:   TimeTo,
				})
				if err != nil {
					return
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
		case "exit":
			return
		}
	}
}
