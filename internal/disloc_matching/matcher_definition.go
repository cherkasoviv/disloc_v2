package disloc_matching

import (
	"errors"
	"github.com/cherkasoviv/go_disl/internal/disloc_storage"
	"github.com/cherkasoviv/go_disl/internal/equipment_event"
	"github.com/cherkasoviv/go_disl/internal/equipment_shipment"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type DislocMatcher struct {
	storage          *disloc_storage.MongoStorage
	goroutinesNumber int
}

func Initialize(storage *disloc_storage.MongoStorage, goroutinesNumber int) *DislocMatcher {
	return &DislocMatcher{
		storage:          storage,
		goroutinesNumber: goroutinesNumber}
}

func (matcher *DislocMatcher) Start() {
	newContainersToMatch, err := matcher.storage.FindNewContainers(context.Background())
	if err != nil {
		return
	}
	for _, container := range newContainersToMatch {
		containerNumber := container.(string)
		equipmentEvents := matcher.storage.FindEquipmentsByContainerNumber(containerNumber, context.Background())
		needToFindShipment := true
		var es *equipment_shipment.EquipmentShipment
		var matchedEvents []*equipment_event.EquipmentEvent
		matchedEvents = make([]*equipment_event.EquipmentEvent, 0)
		for id, ee := range equipmentEvents {
			if ee.GetStatus() == startShipmentEventStatusID {
				es = equipment_shipment.NewEquipmentShipment(
					ee.GetContainerNumber(),
					ee.GetWagonNumber(),
					ee.GetOrderID(),
					ee.GetDateTime())

				needToFindShipment = false
				matchedEvents = append(matchedEvents, ee)
				ee.SetObjectID(id)
				err = matcher.storage.WriteEquipmentEvent(ee, "matched", context.Background())
				if err != nil {
					return
				}

			}
			if needToFindShipment {
				es, err = matcher.storage.FindShipmentByContainer(containerNumber, context.Background())
				if errors.Is(err, mongo.ErrNoDocuments) {
					ee.SetObjectID(id)
					err := matcher.storage.WriteEquipmentEvent(ee, "unmatched", context.Background())
					if err != nil {
						return
					}
				} else if err != nil {
					return
				}

			}
			if ee.GetStatus() != startShipmentEventStatusID {
				if es != nil {
					matchedEvents = append(matchedEvents, ee)
					if ee.GetOrderID() != 0 && es.GetOrderID() == 0 {
						es.SetOrderID(ee.GetOrderID())
					}
					if ee.GetWagonNumber() != "" && es.GetWagonNumber() == "" {
						es.SetWagonNumber(ee.GetWagonNumber())
					}
					ee.SetObjectID(id)
					err = matcher.storage.WriteEquipmentEvent(ee, "matched", context.Background())
					if err != nil {
						return
					}

					if ee.GetStatus() == endShipmentEventStatusID {
						es.SetDateEnd(ee.GetDateTime())
						err := matcher.storage.WriteEquipmentShipment(es, context.Background())
						if err != nil {
							return
						}
						for _, me := range matchedEvents {
							matcher.storage.WriteMatchedEvent(me, es, context.Background())
						}
						matchedEvents = make([]*equipment_event.EquipmentEvent, 0)
						needToFindShipment = true
					}

				} else {
					ee.SetObjectID(id)
					err = matcher.storage.WriteEquipmentEvent(ee, "unmatched", context.Background())
					if err != nil {
						return
					}
				}
			}

		}
		if es != nil {
			err := matcher.storage.WriteEquipmentShipment(es, context.Background())
			if err != nil {
				return
			}
		}
	}
}
