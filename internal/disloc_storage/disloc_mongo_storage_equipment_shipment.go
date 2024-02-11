package disloc_storage

import (
	"context"
	"fmt"
	"github.com/cherkasoviv/go_disl/internal/equipment_shipment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mongoStorage *MongoStorage) WriteEquipmentShipment(equipmentShipment *equipment_shipment.EquipmentShipment) error {
	dateStart, dateEnd := equipmentShipment.GetPeriod()
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("equipment_shipments")
	esBSON := bson.D{
		{"order_id", equipmentShipment.GetOrderID()},
		{"container_number", equipmentShipment.GetContainerNumber()},
		{"wagon_number", equipmentShipment.GetWagonNumber()},
		{"date_start", &dateStart},
		{"date_end", &dateEnd},
	}
	if equipmentShipment.GetObjectID().IsZero() {
		_, err := collection.InsertOne(context.TODO(), esBSON)
		if err != nil {
			return err
		}
	} else {
		_, err := collection.UpdateByID(context.TODO(), equipmentShipment.GetObjectID(), bson.D{{"$set", esBSON}})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (mongoStorage *MongoStorage) FindShipmentByContainer(containerNumber string) (*equipment_shipment.EquipmentShipment, error) {
	db := mongoStorage.client.Database("dislocation")
	collection := db.Collection("equipment_shipments")
	opts := options.FindOne().SetSort(bson.D{{"date_start", -1}})
	result := collection.FindOne(context.TODO(), bson.D{{"container_number", containerNumber}, {"date_end", nil}}, opts)

	var resultBson bson.M
	err := result.Decode(&resultBson)
	if err != nil {
		return nil, err
	}
	es := equipment_shipment.NewEquipmentShipment(
		resultBson["container_number"].(string),
		resultBson["wagon_number"].(string),
		int(resultBson["order_id"].(int32)),
		resultBson["date_start"].(primitive.DateTime).Time(),
	)
	es.SetObjectID(resultBson["_id"].(primitive.ObjectID))
	fmt.Println(es)
	return es, nil

}
