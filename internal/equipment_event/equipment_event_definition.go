package equipment_event

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EquipmentEvent struct {
	objectID        primitive.ObjectID
	statusID        string
	datetime        time.Time
	orderID         int
	containerNumber string
	wagonNumber     string
	stationID       string
}

func New(statusID string, datetime time.Time, orderID int, container string, wagon string, station string) *EquipmentEvent {
	ee := EquipmentEvent{
		statusID:        statusID,
		datetime:        datetime,
		orderID:         orderID,
		containerNumber: container,
		wagonNumber:     wagon,
		stationID:       station,
	}
	return &ee
}

func (ee *EquipmentEvent) GetStatus() string {
	return ee.statusID
}

func (ee *EquipmentEvent) GetDateTime() time.Time {
	return ee.datetime
}

func (ee *EquipmentEvent) GetOrderID() int {
	return ee.orderID
}

func (ee *EquipmentEvent) GetContainerNumber() string {
	return ee.containerNumber

}

func (ee *EquipmentEvent) GetWagonNumber() string {
	return ee.wagonNumber

}

func (ee *EquipmentEvent) GetStationID() string {
	return ee.stationID
}

func (ee *EquipmentEvent) GetObjectID() primitive.ObjectID {
	return ee.objectID
}

func (ee *EquipmentEvent) SetObjectID(id primitive.ObjectID) {
	ee.objectID = id
}
