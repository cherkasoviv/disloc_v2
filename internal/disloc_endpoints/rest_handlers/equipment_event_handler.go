package rest_handlers

import (
	"encoding/json"
	"github.com/cherkasoviv/go_disl/internal/equipment_event"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type eventSaver interface {
	WriteEquipmentEvent(ee *equipment_event.EquipmentEvent, processingStatus string) error
	FindEquipmentsByContainerNumber(containerNumber string) map[primitive.ObjectID]*equipment_event.EquipmentEvent
}

type EquipmentEventHandler struct {
	storage eventSaver
}

type requestForEquipmentEventPOST struct {
	StatusID        string    `json:"status_id"`
	Datetime        time.Time `json:"datetime"`
	OrderID         int       `json:"order_id"`
	ContainerNumber string    `json:"container_number"`
	WagonNumber     string    `json:"wagon_number"`
	StationID       string    `json:"station_id"`
}

type requestForEquipmentEventsByContainerNumberGET struct {
	ContainerNumber string `json:"container_number"`
}

func NewEquipmentEventHandler(str eventSaver) *EquipmentEventHandler {
	return &EquipmentEventHandler{storage: str}
}

func (eeHandler *EquipmentEventHandler) CreateEventFromREST() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			return
		}
		var req requestForEquipmentEventPOST
		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		equipmentEvent := equipment_event.New(
			req.StatusID,
			req.Datetime,
			req.OrderID,
			req.ContainerNumber,
			req.WagonNumber,
			req.StationID,
		)
		err = eeHandler.storage.WriteEquipmentEvent(equipmentEvent, "new")
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)

	}

}

func (eeHandler *EquipmentEventHandler) GetEventsByContainerNumber() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			return
		}
		var req requestForEquipmentEventsByContainerNumberGET
		err := render.DecodeJSON(request.Body, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		containerNumber := req.ContainerNumber

		events := eeHandler.storage.FindEquipmentsByContainerNumber(containerNumber)
		eventsJson, err := json.Marshal(events)

		_, err = writer.Write(eventsJson)
		if err != nil {
			return
		}

	}

}
