package disloc_storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	connectionString string
	client           *mongo.Client
}

// TODO add conn string from config
func InitializeMongoStorage(uri string) (*MongoStorage, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ms := MongoStorage{
		connectionString: "",
		client:           client,
	}
	return &ms, nil
}

//TODO healthcheck

func (mongoStorage *MongoStorage) HealtCheck() error {
	return nil
}
