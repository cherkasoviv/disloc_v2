package test

import (
	"context"
	"fmt"
	rabbit_consumers "github.com/cherkasoviv/go_disl/internal/disloc_endpoints/disloc_consumers"
	"github.com/cherkasoviv/go_disl/internal/disloc_endpoints/disloc_producers"
	"github.com/cherkasoviv/go_disl/internal/disloc_matching"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
	"time"
)

var mongoPort string
var rabbitPort string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	resourceMongo, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env:        []string{},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		dbClient, err := mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://localhost:%s", resourceMongo.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		mongoPort = resourceMongo.GetPort("27017/tcp")
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resourceRabbit, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "rabbitmq",

		Env: []string{},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	rabbitPort = resourceRabbit.GetPort("5672/tcp")
	fmt.Println(rabbitPort)
	fmt.Println(mongoPort)
	time.Sleep(60 * time.Second)
	code := m.Run()
	// When you're done, kill and remove the container
	if err = pool.Purge(resourceMongo); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	if err = pool.Purge(resourceRabbit); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client

	os.Exit(code)
}

func TestConsumer(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:" + rabbitPort + "/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("error getting a channel: %s", err)
	}

	storage, err := disloc_storage.InitializeMongoStorage("mongodb://localhost:" + mongoPort + "/?readPreference=primary&directConnection=true&ssl=false")
	if err != nil {
		return
	}
	rabbit_consumers.NewEquipmentEventConsumer(
		"amqp://guest:guest@localhost:"+rabbitPort,
		"ESB",
		"topic",
		"ee_queue",
		"equipment_event",
		"ee-consumer",
		storage,
	)

	err = channel.PublishWithContext(context.TODO(),
		"ESB",             // exchange
		"equipment_event", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("{\n \"status_id\": \"0604\",\n  \"datetime\":  \"2024-01-04T05:02:02.000Z\",\n  \"order_id\" : 31352124,\n  \"container_number\": \"TTTU4488225\",\n  \"wagon_number\": \"98015837\",\n  \"station_id\": \"002003\"\n}"),
		})
	if err != nil {
		log.Fatalf("producer: error in publish: %s", err)

	}
	time.Sleep(20 * time.Second)
	i := len(storage.FindEquipmentsByContainerNumber("TTTU4488225"))
	t.Run("simple count", func(t *testing.T) {
		if i != 1 {
			t.Errorf("Different number")
		}
	})
	err = channel.PublishWithContext(context.TODO(),
		"ESB",             // exchange
		"equipment_event", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("{\n  \"status_id\": \"1821\",\n  \"datetime\":  \"2024-01-04T05:02:02.000Z\",\n  \"container_number\": \"TTTU4488225\",\n \"wagon_number\": \"98015837\",\n  \"station_id\": \"002003\"\n}"),
		})
	if err != nil {
		log.Fatalf("producer: error in publish: %s", err)
	}
	time.Sleep(20 * time.Second)
	i = len(storage.FindEquipmentsByContainerNumber("TTTU4488225"))
	t.Run("simple second count", func(t *testing.T) {
		if i != 2 {
			t.Errorf("Different number")
		}
	})
	err = channel.PublishWithContext(context.TODO(),
		"ESB",             // exchange
		"equipment_event", // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("{\n  \"status_id\": \"1401\",\n  \"datetime\":  \"2024-01-04T05:02:02.000Z\",\n  \"container_number\": \"TTTU4488225\",\n  \"wagon_number\": \"98015837\",\n  \"station_id\": \"002003\"\n}"),
		})
	if err != nil {
		log.Fatalf("producer: error in publish: %s", err)
	}
	time.Sleep(20 * time.Second)
	i = len(storage.FindEquipmentsByContainerNumber("TTTU4488225"))
	t.Run("simple second count", func(t *testing.T) {
		if i != 3 {
			t.Errorf("Different number")
		}
	})

	matcher := disloc_matching.Initialize(storage, 1)
	matcher.Start()

	eventSender := disloc_producers.NewEventProducer("amqp://guest:guest@localhost:" + rabbitPort + "/")
	eventSender.Publish(storage)
	time.Sleep(20 * time.Second)
	var next bool
	next = true
	evnts := []amqp.Delivery{}
	for next {
		msg, ok, err := channel.Get("matched_events", true)
		if err != nil {
			t.Errorf("Publisher error")
		}
		next = ok
		if next {
			evnts = append(evnts, msg)
		}

	}
	mi := len(evnts)
	fmt.Println(mi)
	if mi != 3 {
		t.Errorf("Not enough matched events")
	}
}
