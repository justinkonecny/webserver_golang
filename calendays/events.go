package calendays

import (
	"../common"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateCalendaysRequest(w, r)
	if !auth {
		common.ErrorUnauthorized(w, r)
		return
	}

	userID := values[KeyUserID].(uint)
	username := values[KeyUsername].(string)

	switch r.Method {
	case http.MethodGet:
		handleEventsGet(w, userID)
	case http.MethodPost:
		handleEventsPost(w, r, username)
	case http.MethodPut:
		fmt.Println("PUT /events")
	case http.MethodDelete:
		fmt.Println("DELETE /events")
	default:
		common.ErrorMethodNotAllowed(w, r)
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

	common.WriteJsonResponse(w, eventsResponse)
}

func handleEventsPost(w http.ResponseWriter, r *http.Request, username string) {
	fmt.Println("POST /events")
	var eventDTO DTOEvent
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&eventDTO)
	if err != nil {
		common.ErrorBadRequest(w, r, err)
		return
	}

	event := Event{
		Name:      eventDTO.Name,
		StartDate: eventDTO.StartDate,
		EndDate:   eventDTO.EndDate,
		Location:  eventDTO.Location,
		Message:   eventDTO.Message,
		NetworkId: eventDTO.NetworkId,
	}
	DB.Create(&event)
	common.WriteJsonResponseWithStatus(w, ConvertEvent(event), http.StatusCreated)
}
