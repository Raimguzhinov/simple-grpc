package user

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	eventmanager "github.com/Raimguzhinov/simple-grpc/pkg/delivery/grpc"
)

func notifyHandler(event *eventmanager.Event, exchange string, routingKey string, queueName string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("Unable to connect to RabbitMQ. Error: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Unable to open a channel. Error: %s", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(
		exchange, // exchange name
		"direct", // type
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		fmt.Printf("Unable to declare exchange. Error: %s", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		fmt.Printf("Unable to declare queue. Error: %s", err)
	}

	if err := ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		fmt.Printf("Unable to bind queue. Error: %s", err)
	}

	messages, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		fmt.Printf("Unable to consume messages. Error: %s", err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			err := proto.Unmarshal(message.Body, event)
			if err != nil {
				panic(err)
			}
			fmt.Println("Notification!")
			t := time.UnixMilli(event.Time).Local().Format(time.DateTime)
			fmt.Printf("Event {\n  senderId: %d\n  eventId: %d\n  time: %s\n  name: '%s'\n}\n> ", event.SenderId, event.EventId, t, event.Name)
		}
	}()

	<-forever
}
