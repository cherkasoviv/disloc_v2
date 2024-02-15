package main

import (
	"context"
	"github.com/cherkasoviv/go_disl/internal/disloc_matching"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
)

func main() {

	storage, err := disloc_storage.InitializeMongoStorage("mongodb://localhost:27017/?readPreference=primary&directConnection=true&ssl=false", context.Background())
	if err != nil {
		return
	}

	matcher := disloc_matching.Initialize(storage, 1)
	matcher.Start()

}
