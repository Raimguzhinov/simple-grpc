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

func publish(pubChan chan models.Events) {
	const exchange = "event.ex"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go func() {
		var conn *amqp.Connection
		var ch *amqp.Channel

		redial := func() {
			var err error
			conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
			if err != nil {
				log.Printf("Unable to connect to RabbitMQ. Error: %s", err)
			} else {
				log.Println("Connected to RabbitMQ")
			}
			ch, err = conn.Channel()
			if err != nil {
				log.Printf("Unable to open a channel. Error: %s", err)
			}
			if err := ch.ExchangeDeclare(
				exchange, // exchange name
				"direct", // type
				true,     // durable
				false,    // delete when unused
				false,    // exclusive
				false,    // no-wait
				nil,      // arguments
			); err != nil {
				log.Printf("Unable to declare exchange. Error: %s", err)
			}
		}

		redial()

		go func() {
			errChan := make(chan *amqp.Error)
			conn.NotifyClose(errChan)
			for err := range errChan {
				log.Printf("Connection to RabbitMQ closed. Reconnecting... Reason: %v", err)
				redial()
			}
		}()

		for {
			event := <-pubChan

			t1 := time.Now().UTC()
			t2 := time.UnixMilli(event.Time).UTC()
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
			err = ch.PublishWithContext(ctx, exchange, strconv.Itoa(int(event.SenderID)), false, false, amqp.Publishing{
				Body: body,
			})
			if err != nil {
				log.Fatalf("Unable to publish event. Error: %s", err)
			}
			log.Println("[x] Sent:", body)
		}
	}()
}
