package main

import (
	"fmt"
	rabbit_consumers "github.com/cherkasoviv/go_disl/internal/disloc_endpoints/disloc_consumers"
	"github.com/cherkasoviv/go_disl/internal/disloc_endpoints/rest_handlers"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	fmt.Println("test")
	storage, err := disloc_storage.InitializeMongoStorage()
	if err != nil {
		return
	}

	eventHandler := rest_handlers.NewEquipmentEventHandler(storage)

	r := chi.NewRouter()
	r.Route("/equipment_event", func(r chi.Router) {
		r.Post("/", eventHandler.CreateEventFromREST())
		r.Get("/", eventHandler.GetEventsByContainerNumber())
	})

	eeConsumer, _ := rabbit_consumers.NewEquipmentEventConsumer(
		"amqp://guest:guest@localhost",
		"ESB",
		"topic",
		"ee_queue",
		"equipment_event",
		"ee-consumer",
		storage,
	)

	defer eeConsumer.Shutdown()

	err = http.ListenAndServe("localhost:8080", r)

	if err != nil {
		panic(err)
	}
}
