package user

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func RunEventsClient(client eventmanager.EventsClient, senderID *int64) {
	var (
		procedureName string
		eventID       int64
		eventName     string
		timeCer       string
		timeFrom      string
		timeTo        string
	)

	exchange := "event.ex"
	routingKey := strconv.Itoa(int(*senderID))
	queueName := strconv.Itoa(int(*senderID))

	go notifyHandler(&eventmanager.Event{}, exchange, routingKey, queueName)

	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		fmt.Print("> ")
		fmt.Scan(&procedureName)
		switch procedureName {
		case "MakeEvent":
			eventMaker(client, senderID, timeCer, eventName)
		case "GetEvent":
			eventGetter(client, senderID, eventID)
		case "DeleteEvent":
			eventDeleter(client, senderID, eventID)
		case "GetEvents":
			eventsGetter(client, senderID, timeFrom, timeTo)
		case "exit":
			return
		}
	}
}

func eventMaker(client eventmanager.EventsClient, senderID *int64, timeCer string, eventName string) {
	fmt.Print("Enter <time> as format 2006-01-02(15:04) and <event_name>: ")
	fmt.Scan(&timeCer, &eventName)
	datetime, err := time.ParseInLocation("2006-01-02(15:04)", timeCer, time.Local)
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := client.MakeEvent(context.Background(), &eventmanager.MakeEventRequest{
			SenderId: *senderID,
			Time:     datetime.UnixMilli(),
			Name:     eventName,
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Created {", res, "}")
	}
}

func eventGetter(client eventmanager.EventsClient, senderID *int64, eventID int64) {
	fmt.Print("Enter <event_id>: ")
	if _, err := fmt.Scan(&eventID); err != nil {
		fmt.Println(err)
	} else {
		res, err := client.GetEvent(context.Background(), &eventmanager.GetEventRequest{
			SenderId: *senderID,
			EventId:  eventID,
		})
		if err != nil {
			parts := strings.Split(err.Error(), "desc = ")
			fmt.Println(parts[1])
		} else {
			t := time.UnixMilli(res.Time).Local().Format("2006-01-02(15:04)")
			fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n", res.SenderId, res.EventId, t, res.Name)
		}
	}
}

func eventDeleter(client eventmanager.EventsClient, senderID *int64, eventID int64) {
	fmt.Print("Enter <event_id>: ")
	if _, err := fmt.Scan(&eventID); err != nil {
		fmt.Println(err)
	} else {
		res, err := client.DeleteEvent(context.Background(), &eventmanager.DeleteEventRequest{
			SenderId: *senderID,
			EventId:  eventID,
		})
		if err != nil {
			parts := strings.Split(err.Error(), "desc = ")
			fmt.Println(parts[1])
		} else {
			fmt.Println("Deleted {", res, "}")
		}
	}
}

func eventsGetter(client eventmanager.EventsClient, senderID *int64, timeFrom string, timeTo string) {
	fmt.Print("Enter <from_time> <to_time> as format 2006-01-02(15:04): ")
	if _, err := fmt.Scan(&timeFrom, &timeTo); err != nil {
		fmt.Println(err)
	} else {
		datetimeFrom, err1 := time.ParseInLocation("2006-01-02(15:04)", timeFrom, time.Local)
		datetimeTo, err2 := time.ParseInLocation("2006-01-02(15:04)", timeTo, time.Local)

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
						t := time.UnixMilli(res.Time).Local().Format("2006-01-02(15:04)")
						fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n", res.SenderId, res.EventId, t, res.Name)
					}
				}
			}
		}
	}
}
