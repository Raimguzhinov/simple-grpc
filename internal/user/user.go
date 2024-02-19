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
		dateCer       string
		dateFrom      string
		dateTo        string
		timeCer       string
		timeFrom      string
		timeTo        string
	)

	routingKey := strconv.Itoa(int(*senderID))
	queueName := strconv.Itoa(int(*senderID))

	go notifyHandler(&eventmanager.Event{}, routingKey, queueName)

	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		fmt.Print("> ")
		if _, err := fmt.Scan(&procedureName); err != nil {
			fmt.Println(err)
		}
		switch procedureName {
		case "MakeEvent":
			eventMaker(client, senderID, dateCer, timeCer, eventName)
		case "GetEvent":
			eventGetter(client, senderID, eventID)
		case "DeleteEvent":
			eventDeleter(client, senderID, eventID)
		case "GetEvents":
			eventsGetter(client, senderID, dateFrom, timeFrom, dateTo, timeTo)
		case "exit":
			return
		default:
			fmt.Println("Bad Procedure Name")
		}
	}
}

func eventMaker(client eventmanager.EventsClient, senderID *int64, dateCer string, timeCer string, eventName string) {
	fmt.Print("Enter <date> <time> <event_name>: ")
	if _, err := fmt.Scan(&dateCer, &timeCer, &eventName); err != nil {
		fmt.Println(err)
	}
	locDateTime, err := time.ParseInLocation(time.DateTime, dateCer+" "+timeCer, time.Local)
	dateTime := locDateTime.UTC()
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := client.MakeEvent(context.Background(), &eventmanager.MakeEventRequest{
			SenderId: *senderID,
			Time:     dateTime.UnixMilli(),
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
			t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
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

func eventsGetter(client eventmanager.EventsClient, senderID *int64, dateFrom string, timeFrom string, dateTo string, timeTo string) {
	fmt.Print("Enter <from_date> <from_time> <to_date> <to_time>: ")
	if _, err := fmt.Scan(&dateFrom, &timeFrom, &dateTo, &timeTo); err != nil {
		fmt.Println(err)
	} else {
		locDateTimeFrom, err1 := time.ParseInLocation(time.DateTime, dateFrom+" "+timeFrom, time.Local)
		locDateTimeTo, err2 := time.ParseInLocation(time.DateTime, dateTo+" "+timeTo, time.Local)
		dateTimeFrom := locDateTimeFrom.UTC()
		dateTimeTo := locDateTimeTo.UTC()

		if err1 != nil || err2 != nil {
			fmt.Println(err1, err2)
		} else {
			stream, err := client.GetEvents(context.Background(), &eventmanager.GetEventsRequest{
				SenderId: *senderID,
				FromTime: dateTimeFrom.UnixMilli(),
				ToTime:   dateTimeTo.UnixMilli(),
			})
			if err != nil {
				fmt.Println(err)
			} else {
				for i := 0; ; i++ {
					res, err := stream.Recv()
					if err == io.EOF {
						if i == 0 {
							fmt.Println("Not found")
						}
						break
					}
					t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
					fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n", res.SenderId, res.EventId, t, res.Name)
				}
			}
		}
	}
}
