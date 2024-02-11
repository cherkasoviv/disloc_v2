package disloc_storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	connString string
	client     *mongo.Client
}

// TODO add conn string from config
func InitializeMongoStorage() (*MongoStorage, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&directConnection=true&ssl=false"))
	if err != nil {
		return nil, err
	}
	ms := MongoStorage{
		connString: "",
		client:     client,
	}
	return &ms, nil
}

//TODO healthcheck

func (mongoStorage *MongoStorage) HealtCheck() error {
	return nil
}
