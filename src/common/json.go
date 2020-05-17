package common

import "time"

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

func ConvertNetworkList(networks []Network) []DTONetwork {
	var networksDTO []DTONetwork
	for _, n := range networks {
		out := DTONetwork{
			ID:      n.ID,
			Name:    n.Name,
			OwnerId: n.UserId,
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
