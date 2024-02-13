package service

import (
	"context"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	"github.com/Raimguzhinov/simple-grpc/internal/models"
)

func publish(event models.Events) {
	exchange := "event.ex"
	routingKey := strconv.Itoa(int(event.SenderID))
	queueName := strconv.Itoa(int(event.SenderID))

	go func() {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Fatalf("Unable to connect to RabbitMQ. Error: %s", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Unable to open a channel. Error: %s", err)
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
			log.Fatalf("Unable to declare exchange. Error: %s", err)
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
			log.Fatalf("Unable to declare queue. Error: %s", err)
		}

		if err := ch.QueueBind(
			q.Name,     // queue name
			routingKey, // routing key
			exchange,   // exchange
			false,      // no-wait
			nil,        // arguments
		); err != nil {
			log.Fatalf("Unable to bind queue. Error: %s", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		t1 := time.Now().Local()
		t2 := time.UnixMilli(event.Time)
		timeDuration := t2.Sub(t1)
		if timeDuration < 0 {
			return
		}
		timer := time.NewTimer(timeDuration)
		<-timer.C
		body, err := proto.Marshal(accumulateEvent(event))
		if err != nil {
			log.Fatalf("Unable to marshal event. Error: %s", err)
		}

		err = ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
			Body: body,
		})
		if err != nil {
			log.Fatalf("Unable to publish event. Error: %s", err)
		}

		log.Println("[x] Sent:", body)
	}()
}
