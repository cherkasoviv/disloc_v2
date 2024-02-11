package main

import (
	"github.com/cherkasoviv/go_disl/internal/disloc_matching"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
)

func main() {

	storage, err := disloc_storage.InitializeMongoStorage()
	if err != nil {
		return
	}

	matcher := disloc_matching.Initialize(storage, 1)
	matcher.Start()
	storage.FindShipmentByContainer("TTTU4488223")
}
