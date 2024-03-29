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
	"github.com/google/uuid"

	eventctrl "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

type ClientManager interface {
	EventMaker(senderID int64, dateCer string, timeCer string, eventName string)
	EventGetter(senderID int64, eventID []byte)
	EventDeleter(senderID int64, eventID []byte)
	EventsGetter(senderID int64, dateFrom string, timeFrom string, dateTo string, timeTo string)
}

type User struct {
	eventctrl.EventsClient
}

func NewUser(conn *eventctrl.EventsClient) *User {
	return &User{
		EventsClient: *conn,
	}
}

func RunEventsClient(client eventctrl.EventsClient, senderID int64) {
	var eventID []byte
	var eventName,
		dateCer,
		dateFrom,
		dateTo,
		timeCer,
		timeFrom,
		timeTo string

	u := NewUser(&client)

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
			u.EventMaker(senderID, dateCer, timeCer, eventName)
		case "GetEvent":
			u.EventGetter(senderID, eventID)
		case "DeleteEvent":
			u.EventDeleter(senderID, eventID)
		case "GetEvents":
			u.EventsGetter(senderID, dateFrom, timeFrom, dateTo, timeTo)
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

func (u *User) EventMaker(
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
		res, err := u.MakeEvent(context.Background(), &eventctrl.MakeEventRequest{
			SenderId: senderID,
			Time:     dateTime.UnixMilli(),
			Name:     eventName,
		})
		if err != nil {
			fmt.Println(err)
		}
		eventID, err := uuid.FromBytes(res.EventId)
		if err != nil {
			return
		}
		fmt.Println("Created {", eventID, "}")
	}
}

func (u *User) EventGetter(senderID int64, eventID []byte) {
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
		res, err := u.GetEvent(context.Background(), &eventctrl.GetEventRequest{
			SenderId: senderID,
			EventId:  eventID,
		})
		if err != nil {
			parts := strings.Split(err.Error(), "desc = ")
			fmt.Println(parts[1])
		} else {
			eventID, err := uuid.FromBytes(res.EventId)
			if err != nil {
				return
			}
			t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
			fmt.Printf("Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n", res.SenderId, eventID, t, res.Name)
		}
	}
}

func (u *User) EventDeleter(senderID int64, eventID []byte) {
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
		res, err := u.DeleteEvent(context.Background(), &eventctrl.DeleteEventRequest{
			SenderId: senderID,
			EventId:  eventID,
		})
		if err != nil {
			parts := strings.Split(err.Error(), "desc = ")
			fmt.Println(parts[1])
		} else {
			eventID, err := uuid.FromBytes(res.EventId)
			if err != nil {
				return
			}
			fmt.Println("Deleted {", eventID, "}")
		}
	}
}

func (u *User) EventsGetter(
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
			stream, err := u.GetEvents(context.Background(), &eventctrl.GetEventsRequest{
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
					eventID, err := uuid.FromBytes(res.EventId)
					if err != nil {
						return
					}
					t := time.UnixMilli(res.Time).Local().Format(time.DateTime)
					fmt.Printf("Event {\n  senderId: %d\n  eventId: %s\n  time: %s\n  name: '%s'\n}\n", res.SenderId, eventID, t, res.Name)
				}
			}
		}
	}
}
