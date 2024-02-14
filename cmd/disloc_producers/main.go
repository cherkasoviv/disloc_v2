package main

import (
	"github.com/cherkasoviv/go_disl/internal/disloc_endpoints/disloc_producers"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
)

func main() {

	storage, err := disloc_storage.InitializeMongoStorage("mongodb://localhost:27017/?readPreference=primary&directConnection=true&ssl=false")
	if err != nil {
		return
	}
	eventSender := disloc_producers.NewEventProducer("amqp://guest:guest@localhost:5672/")
	eventSender.Publish(storage)

}
