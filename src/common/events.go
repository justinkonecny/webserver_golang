package common

import (
	"fmt"
	"net/http"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateRequest(r)
	if !auth {
		ErrorUnauthorized(w)
		return
	}

	userID := values[KeyUserID].(uint)

	switch r.Method {
	case http.MethodGet:
		handleEventsGet(w, userID)
	case http.MethodPost:
		fmt.Println("POST /events")
	case http.MethodPut:
		fmt.Println("PUT /events")
	case http.MethodDelete:
		fmt.Println("DELETE /events")
	default:
		ErrorMethodNotAllowed(w)
	}
}

func handleEventsGet(w http.ResponseWriter, userID uint) {
	fmt.Println("GET /events")
	var eventsResponse []DTOEvent
	var networkUsers []NetworkUser
	DB.Where(&NetworkUser{UserId: userID}).Find(&networkUsers)

	for _, networkUser := range networkUsers {
		networkID := networkUser.NetworkId
		var events []Event
		DB.Where(&Event{NetworkId: networkID}).Find(&events)
		eventsDTO := ConvertEventList(events, networkID)
		eventsResponse = append(eventsResponse, eventsDTO...)
	}

	WriteJsonResponse(w, eventsResponse)
}
