package disloc_consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
	"github.com/cherkasoviv/go_disl/internal/equipment_event"
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
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	tag     string
}

func NewEquipmentEventConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string, storage *disloc_storage.MongoStorage) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		true,
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done, storage)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error, storage *disloc_storage.MongoStorage) {
	cleanup := func() {
		log.Printf("handle: deliveries channel closed")
		done <- nil
	}

	defer cleanup()

	for d := range deliveries {
		equipmentEventMessage := amqpEquipmentEvent{}
		err := json.Unmarshal(d.Body, &equipmentEventMessage)
		if err != nil {
			return
		}
		equipmentEvent := equipment_event.New(
			equipmentEventMessage.StatusID,
			equipmentEventMessage.Datetime,
			equipmentEventMessage.OrderID,
			equipmentEventMessage.ContainerNumber,
			equipmentEventMessage.WagonNumber,
			equipmentEventMessage.StationID,
		)
		err = storage.WriteEquipmentEvent(equipmentEvent, "new", context.Background())
		if err != nil {
			return
		}

	}
}
