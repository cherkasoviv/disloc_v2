package equipment_shipment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EquipmentShipment struct {
	objectID        primitive.ObjectID
	containerNumber string
	wagonNumber     string
	orderID         int
	dateStart       *time.Time
	dateEnd         *time.Time
}

func (receiver EquipmentShipment) name() {

}

func NewEquipmentShipment(containerNumber string, wagonNumber string, orderID int, dateStart time.Time) *EquipmentShipment {
	return &EquipmentShipment{
		containerNumber: containerNumber,
		wagonNumber:     wagonNumber,
		orderID:         orderID,
		dateStart:       &dateStart,
		dateEnd:         nil,
	}
}

func (equipmentShipment *EquipmentShipment) GetContainerNumber() string {
	return equipmentShipment.containerNumber
}

func (equipmentShipment *EquipmentShipment) GetWagonNumber() string {
	return equipmentShipment.wagonNumber
}

func (equipmentShipment *EquipmentShipment) GetOrderID() int {
	return equipmentShipment.orderID
}

func (equipmentShipment *EquipmentShipment) GetPeriod() (*time.Time, *time.Time) {
	return equipmentShipment.dateStart, equipmentShipment.dateEnd
}

func (equipmentShipment *EquipmentShipment) SetDateEnd(dateEnd time.Time) {
	equipmentShipment.dateEnd = &dateEnd

}

func (equipmentShipment *EquipmentShipment) SetObjectID(objectID primitive.ObjectID) {
	equipmentShipment.objectID = objectID

}

func (equipmentShipment *EquipmentShipment) SetOrderID(id int) {
	equipmentShipment.orderID = id
}

func (equipmentShipment *EquipmentShipment) SetWagonNumber(number string) {
	equipmentShipment.wagonNumber = number
}

func (equipmentShipment *EquipmentShipment) GetObjectID() primitive.ObjectID {
	return equipmentShipment.objectID
}
