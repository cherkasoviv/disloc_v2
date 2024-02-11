package equipment_event

import (
	"reflect"
	"testing"
	"time"
)

func TestEquipmentEvent_GetStatus(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get_status_id",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: "1801",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetStatus(); got != tt.want {
				t.Errorf("GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentEvent_GetDateTime(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get_datetime",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetDateTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentEvent_GetOrderID(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get_order_id",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: 31720523,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetOrderID(); got != tt.want {
				t.Errorf("GetOrderID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentEvent_GetContainerNumber(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get_container_number",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: "FFAU4488220",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetContainerNumber(); got != tt.want {
				t.Errorf("GetContainerNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentEvent_GetWagonNumber(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get_wagon_number",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: "98015837",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetWagonNumber(); got != tt.want {
				t.Errorf("GetWagonNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquipmentEvent_GetStationID(t *testing.T) {
	type fields struct {
		statusID        string
		datetime        time.Time
		orderID         int
		containerNumber string
		wagonNumber     string
		stationID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{

		{
			name: "get_station_id",
			fields: struct {
				statusID        string
				datetime        time.Time
				orderID         int
				containerNumber string
				wagonNumber     string
				stationID       string
			}{
				statusID:        "1801",
				datetime:        time.Date(2024, time.February, 2, 2, 2, 2, 2, time.Local),
				orderID:         31720523,
				containerNumber: "FFAU4488220",
				wagonNumber:     "98015837",
				stationID:       "002003",
			},
			want: "002003",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ee := &EquipmentEvent{
				statusID:        tt.fields.statusID,
				datetime:        tt.fields.datetime,
				orderID:         tt.fields.orderID,
				containerNumber: tt.fields.containerNumber,
				wagonNumber:     tt.fields.wagonNumber,
				stationID:       tt.fields.stationID,
			}
			if got := ee.GetStationID(); got != tt.want {
				t.Errorf("GetStationID() = %v, want %v", got, tt.want)
			}
		})
	}
}
