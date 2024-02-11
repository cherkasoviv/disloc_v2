package disloc_storage

import (
	"context"
	"fmt"
	"github.com/cherkasoviv/go_disl/internal/equipment_event"
	"github.com/cherkasoviv/go_disl/internal/equipment_shipment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mongoStorage *MongoStorage) WriteEquipmentEvent(ee *equipment_event.EquipmentEvent, processingStatus string) error {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("incoming_events")
	eeBSON := bson.D{
		{"status_id", ee.GetStatus()},
		{"datetime", ee.GetDateTime()},
		{"order_id", ee.GetOrderID()},
		{"container_number", ee.GetContainerNumber()},
		{"wagon_number", ee.GetWagonNumber()},
		{"station_id", ee.GetStationID()},
		{"status", processingStatus},
	}
	if ee.GetObjectID().IsZero() {
		_, err := collection.InsertOne(context.TODO(), eeBSON)
		if err != nil {
			return err
		}
	} else {
		_, err := collection.UpdateByID(context.TODO(), ee.GetObjectID(), bson.D{{"$set", eeBSON}})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (mongoStorage *MongoStorage) FindNewContainers() ([]interface{}, error) {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("incoming_events")
	distinct, err := collection.Distinct(context.TODO(), "container_number", bson.D{{"status", "new"}})
	if err != nil {
		return nil, err
	}
	return distinct, nil
}

func (mongoStorage *MongoStorage) FindNewWagons() {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("incoming_events")
	distinct, err := collection.Distinct(context.TODO(), "wagon_number", bson.D{{"status", "new"}, {"container_number", ""}})
	if err != nil {
		return
	}
	fmt.Println(distinct)
}

func (mongoStorage *MongoStorage) FindEquipmentsByContainerNumber(containerNumber string) map[primitive.ObjectID]*equipment_event.EquipmentEvent {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("incoming_events")
	opts := options.Find().SetSort(bson.D{{"datetime", 1}})
	result, err := collection.Find(context.TODO(), bson.D{{"container_number", containerNumber}}, opts)
	if err != nil {
		return nil
	}
	var resultBson []bson.M
	err = result.All(context.TODO(), &resultBson)
	if err != nil {
		return nil
	}
	events := map[primitive.ObjectID]*equipment_event.EquipmentEvent{}
	for _, r := range resultBson {

		events[r["_id"].(primitive.ObjectID)] = equipment_event.New(
			r["status_id"].(string),
			r["datetime"].(primitive.DateTime).Time(),
			int(r["order_id"].(int32)),
			r["container_number"].(string),
			r["wagon_number"].(string),
			r["station_id"].(string),
		)
	}

	return events
}

func (mongoStorage *MongoStorage) WriteMatchedEvent(ee *equipment_event.EquipmentEvent, es *equipment_shipment.EquipmentShipment) error {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("outcoming_events")
	eeBSON := bson.D{
		{"status_id", ee.GetStatus()},
		{"datetime", ee.GetDateTime()},
		{"order_id", es.GetOrderID()},
		{"container_number", es.GetContainerNumber()},
		{"wagon_number", es.GetWagonNumber()},
		{"station_id", ee.GetStationID()},
		{"status", "not_sent"},
	}

	_, err := collection.InsertOne(context.TODO(), eeBSON)
	if err != nil {
		return err
	}

	return nil
}

func (mongoStorage *MongoStorage) FindEquipmentToSend() map[primitive.ObjectID]*equipment_event.EquipmentEvent {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("outcoming_events")
	opts := options.Find().SetSort(bson.D{{"datetime", 1}})
	result, err := collection.Find(context.TODO(), bson.D{{"status", "not_sent"}}, opts)
	if err != nil {
		return nil
	}
	var resultBson []bson.M
	err = result.All(context.TODO(), &resultBson)
	if err != nil {
		return nil
	}
	events := map[primitive.ObjectID]*equipment_event.EquipmentEvent{}
	for _, r := range resultBson {

		events[r["_id"].(primitive.ObjectID)] = equipment_event.New(
			r["status_id"].(string),
			r["datetime"].(primitive.DateTime).Time(),
			int(r["order_id"].(int32)),
			r["container_number"].(string),
			r["wagon_number"].(string),
			r["station_id"].(string),
		)
	}

	return events
}

func (mongoStorage *MongoStorage) SetStatusSentForMatchedEvent(id primitive.ObjectID) error {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("outcoming_events")
	_, err := collection.UpdateByID(context.TODO(), id, bson.D{{"$set", bson.D{{"status", "sent"}}}})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
