package eventclt

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/api/protobuf"
)

var (
	ProcedureName string
	EventID       int64
	EventName     string
	Time          string
	TimeFrom      string
	TimeTo        string
)

func RunEventsClient(client eventmanager.EventsClient, senderID *int64) {
	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		fmt.Print("> ")
		fmt.Scan(&ProcedureName)
		switch ProcedureName {
		case "MakeEvent":
			fmt.Print("Enter <time> as format 2006-01-02(15:04) and <event_name>: ")
			fmt.Scan(&Time, &EventName)
			datetime, err := time.Parse("2006-01-02(15:04)", Time)
			if err != nil {
				fmt.Println(err)
			} else {
				res, err := client.MakeEvent(context.Background(), &eventmanager.MakeEventRequest{
					SenderId: *senderID,
					Time:     datetime.UnixMilli(),
					Name:     EventName,
				})
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Created {", res, "}")
			}
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
					t := time.UnixMilli(res.Time).UTC().Format("2006-01-02(15:04)")
					fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n", res.SenderId, res.EventId, t, res.Name)
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
			fmt.Print("Enter <from_time> <to_time> as format 2006-01-02(15:04): ")
			if _, err := fmt.Scan(&TimeFrom, &TimeTo); err != nil {
				fmt.Println(err)
			} else {
				datetimeFrom, err1 := time.Parse("2006-01-02(15:04)", TimeFrom)
				datetimeTo, err2 := time.Parse("2006-01-02(15:04)", TimeTo)
				if err1 != nil || err2 != nil {
					fmt.Println(err1, err2)
				} else {
					stream, err := client.GetEvents(context.Background(), &eventmanager.GetEventsRequest{
						SenderId: *senderID,
						FromTime: datetimeFrom.UnixMilli(),
						ToTime:   datetimeTo.UnixMilli(),
					})
					if err != nil {
						fmt.Println(err)
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
								t := time.UnixMilli(res.Time).UTC().Format("2006-01-02(15:04)")
								fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n", res.SenderId, res.EventId, t, res.Name)
							}
						}
					}
				}
			}
		case "exit":
			return
		}
	}
}
