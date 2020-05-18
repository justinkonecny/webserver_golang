package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type DTONetwork struct {
	ID       uint
	Name     string
	OwnerId  uint
	ColorHex string
}

type DTOEvent struct {
	ID        uint
	StartDate time.Time
	EndDate   time.Time
	Location  string
	Message   string
	NetworkId uint
}

func WriteJsonResponse(w http.ResponseWriter, data interface{}) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		ErrorInternalServerError(w)
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	_, jsonErr := w.Write(jsonData)
	return jsonErr == nil
}

func ConvertNetworkUserList(networkUsers []NetworkUser) []DTONetwork {
	var networksDTO []DTONetwork
	for _, nu := range networkUsers {
		out := DTONetwork{
			ID:       nu.Network.ID,
			Name:     nu.Network.Name,
			OwnerId:  nu.Network.UserId,
			ColorHex: nu.ColorHex,
		}
		networksDTO = append(networksDTO, out)
	}

	return networksDTO
}

func ConvertEventList(events []Event, networkID uint) []DTOEvent {
	var eventsDTO []DTOEvent
	for _, e := range events {
		out := DTOEvent{
			ID:        e.ID,
			StartDate: e.StartDate,
			EndDate:   e.EndDate,
			Location:  e.Location,
			Message:   e.Message,
			NetworkId: networkID,
		}
		eventsDTO = append(eventsDTO, out)
	}
	return eventsDTO
}
