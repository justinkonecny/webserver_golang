package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type DTOUserSignup struct {
	FirstName    string
	LastName     string
	Email        string
	Username     string
	FirebaseUUID string
}

type DTOUser struct {
	ID                   uint
	FirstName            string
	LastName             string
	Email                string
	Username             string
	SubscriptionStatusID uint
}

type DTONetwork struct {
	ID       uint
	Name     string
	OwnerId  uint
	ColorHex string
	Members  []DTOUser
}

type DTOEvent struct {
	ID        uint
	StartDate time.Time
	EndDate   time.Time
	Name      string
	Location  string
	Message   string
	NetworkId uint
}

func WriteJsonResponse(w http.ResponseWriter, data interface{}) bool {
	return WriteJsonResponseWithStatus(w, data, http.StatusOK)
}

func WriteJsonResponseWithStatus(w http.ResponseWriter, data interface{}, status int) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		ErrorInternalServerError(w)
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
			Members:  ConvertUserList(nu.Network.Users),
		}
		networksDTO = append(networksDTO, out)
	}
	return networksDTO
}

func ConvertNetwork(network Network, colorHex string) DTONetwork {
	return DTONetwork{
		ID:       network.ID,
		Name:     network.Name,
		OwnerId:  network.UserId,
		ColorHex: colorHex,
		Members:  ConvertUserList(network.Users),
	}
}

func ConvertEventList(events []Event, networkID uint) []DTOEvent {
	var eventsDTO []DTOEvent
	for _, e := range events {
		out := DTOEvent{
			ID:        e.ID,
			StartDate: e.StartDate,
			EndDate:   e.EndDate,
			Name:      e.Name,
			Location:  e.Location,
			Message:   e.Message,
			NetworkId: networkID,
		}
		eventsDTO = append(eventsDTO, out)
	}
	return eventsDTO
}

func ConvertEvent(event Event) DTOEvent {
	return DTOEvent{
		ID:        event.ID,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
		Name:      event.Name,
		Location:  event.Location,
		Message:   event.Message,
		NetworkId: event.NetworkId,
	}
}

func ConvertUserList(users []User) []DTOUser {
	var usersDTO []DTOUser
	for _, u := range users {
		out := DTOUser{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Username:  u.Username,
		}
		usersDTO = append(usersDTO, out)
	}
	return usersDTO
}

func ConvertUser(user User) DTOUser {
	return DTOUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Username:  user.Username,
	}
}
