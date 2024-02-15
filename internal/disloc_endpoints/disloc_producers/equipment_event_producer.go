package disloc_producers

import (
	"context"
	"encoding/json"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type amqpEquipmentEvent struct {
	StatusID        string    `json:"status_id"`
	Datetime        time.Time `json:"datetime"`
	OrderID         int       `json:"order_id"`
	ContainerNumber string    `json:"container_number"`
	WagonNumber     string    `json:"wagon_number"`
	StationID       string    `json:"station_id"`
}
type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewEventProducer(uri string) *Producer {
	conn, err := amqp.Dial(uri) // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("error getting a channel: %s", err)
	}

	queue, err := channel.QueueDeclare(
		"matched_events", // name of the queue
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // noWait
		nil,              // arguments
	)
	if err == nil {
		log.Printf("producer: declared queue (%q %d messages, %d consumers)",
			queue.Name, queue.Messages, queue.Consumers)
	} else {
		log.Fatalf("producer: Queue Declare: %s", err)
	}
	return &Producer{
		conn:    conn,
		channel: channel,
		queue:   &queue,
	}

}

func (producer *Producer) Publish(storage *disloc_storage.MongoStorage) {
	events := storage.FindEquipmentToSend(context.Background())
	for objectID, event := range events {
		amqpEvent := amqpEquipmentEvent{
			StatusID:        event.GetStatus(),
			Datetime:        event.GetDateTime(),
			OrderID:         event.GetOrderID(),
			ContainerNumber: event.GetContainerNumber(),
			WagonNumber:     event.GetWagonNumber(),
			StationID:       event.GetStationID(),
		}
		jsonEvent, _ := json.Marshal(amqpEvent)
		err := producer.channel.PublishWithContext(context.TODO(),
			"",                  // exchange
			producer.queue.Name, // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        jsonEvent,
			})
		if err != nil {
			log.Fatalf("producer: error in publish: %s", err)

		}
		err = storage.SetStatusSentForMatchedEvent(objectID, context.Background())
		if err != nil {
			return
		}

	}

}
