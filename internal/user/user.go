package user

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"

	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func RunEventsClient(client eventctrl.EventsClient, senderID int64) {
	var (
		eventID   int64
		eventName string
		dateCer   string
		dateFrom  string
		dateTo    string
		timeCer   string
		timeFrom  string
		timeTo    string
	)

	routingKey := strconv.Itoa(int(senderID))
	queueName := strconv.Itoa(int(senderID))

	f := bufio.NewWriter(os.Stdout)
	go notifyHandler(&eventctrl.Event{}, routingKey, queueName, f)

	fmt.Println("Choose procedure: MakeEvent, GetEvent, DeleteEvent, GetEvents")
	for {
		f.Flush()
		sp := selection.New(
			"\nWhat do you pick?",
			[]string{"MakeEvent", "GetEvent", "DeleteEvent", "GetEvents", "exit"},
		)
		sp.PageSize = 4
		choice, err := sp.RunPrompt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		switch choice {
		case "MakeEvent":
			eventMaker(client, senderID, dateCer, timeCer, eventName)
		case "GetEvent":
			eventGetter(client, senderID, eventID)
		case "DeleteEvent":
			eventDeleter(client, senderID, eventID)
		case "GetEvents":
			eventsGetter(client, senderID, dateFrom, timeFrom, dateTo, timeTo)
		case "exit":
			input := confirmation.New("Are you confirming the exit?", confirmation.Yes)
			ready, err := input.RunPrompt()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(1)
			}
			if ready {
				os.Exit(0)
			}
		default:
			fmt.Println("Bad Procedure Name")
		}
	}
}

func eventMaker(
	client eventctrl.EventsClient,
	senderID int64,
	dateCer string,
	timeCer string,
	eventName string,
) {
	input := textinput.New("Enter <date> <time> <event_name>:")
	input.Placeholder = "Args cannot be empty"
	input.Validate = func(input string) error {
		if len(input) < 21 {
			return fmt.Errorf("few arguments")
		}
		return nil
	}
	args, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if _, err = fmt.Sscan(args, &dateCer, &timeCer, &eventName); err != nil {
		fmt.Println(err)
	}
	locDateTime, err := time.ParseInLocation(time.DateTime, dateCer+" "+timeCer, time.Local)
	dateTime := locDateTime.UTC()
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := client.MakeEvent(context.Background(), &eventctrl.MakeEventRequest{
			SenderId: senderID,
			Time:     dateTime.UnixMilli(),
			Name:     eventName,
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Created {", res, "}")
	}
}

func eventGetter(client eventctrl.EventsClient, senderID int64, eventID int64) {
	input := textinput.New("Enter <event_id>:")
	input.Placeholder = "Arg cannot be empty"
	arg, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if _, err := fmt.Sscan(arg, &eventID); err != nil {
		fmt.Println(err)
	} else {
		res, err := client.GetEvent(context.Background(), &eventctrl.GetEventRequest{
			SenderId: senderID,
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

func eventDeleter(client eventctrl.EventsClient, senderID int64, eventID int64) {
	input := textinput.New("Enter <event_id>:")
	input.Placeholder = "Arg cannot be empty"
	arg, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if _, err := fmt.Sscan(arg, &eventID); err != nil {
		fmt.Println(err)
	} else {
		res, err := client.DeleteEvent(context.Background(), &eventctrl.DeleteEventRequest{
			SenderId: senderID,
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

func eventsGetter(
	client eventctrl.EventsClient,
	senderID int64,
	dateFrom string,
	timeFrom string,
	dateTo string,
	timeTo string,
) {
	input := textinput.New("Enter <from_date> <from_time> <to_date> <to_time>:")
	input.Placeholder = "Args cannot be empty"
	input.Validate = func(input string) error {
		if len(input) <= 38 || len(input) > 39 {
			return fmt.Errorf("few arguments")
		}
		return nil
	}
	args, err := input.RunPrompt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	if _, err = fmt.Sscan(args, &dateFrom, &timeFrom, &dateTo, &timeTo); err != nil {
		fmt.Println(err)
	} else {
		locDateTimeFrom, err1 := time.ParseInLocation(time.DateTime, dateFrom+" "+timeFrom, time.Local)
		locDateTimeTo, err2 := time.ParseInLocation(time.DateTime, dateTo+" "+timeTo, time.Local)
		dateTimeFrom := locDateTimeFrom.UTC()
		dateTimeTo := locDateTimeTo.UTC()

		if err1 != nil || err2 != nil {
			fmt.Println(err1, err2)
		} else {
			stream, err := client.GetEvents(context.Background(), &eventctrl.GetEventsRequest{
				SenderId: senderID,
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
