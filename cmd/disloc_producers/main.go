package main

import (
	"github.com/cherkasoviv/go_disl/internal/disloc_endpoints/disloc_producers"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
)

func main() {

	storage, err := disloc_storage.InitializeMongoStorage()
	if err != nil {
		return
	}
	eventSender := disloc_producers.NewEventProducer()
	eventSender.Publish(storage)

}
