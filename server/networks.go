package server

import (
	"fmt"
	"net/http"
)

func HandleNetworks(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateRequest(w, r)
	if !auth {
		ErrorUnauthorized(w)
		return
	}

	userID := values[KeyUserID].(uint)

	switch r.Method {
	case http.MethodGet:
		handleNetworksGet(w, userID)
	case http.MethodPost:
		fmt.Println("POST /networks")
	case http.MethodPut:
		fmt.Println("PUT /networks")
	case http.MethodDelete:
		fmt.Println("DELETE /networks")
	default:
		ErrorMethodNotAllowed(w)
	}
}

func handleNetworksGet(w http.ResponseWriter, userID uint) {
	fmt.Println("GET /networks")
	var networkUsers []NetworkUser
	DB.Preload("Network").Where(&NetworkUser{UserId: userID}).Find(&networkUsers)
	WriteJsonResponse(w, ConvertNetworkUserList(networkUsers))
}
