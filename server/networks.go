package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HandleNetworks(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateRequest(w, r)
	if !auth {
		ErrorUnauthorized(w, r)
		return
	}

	userID := values[KeyUserID].(uint)
	username := values[KeyUsername].(string)

	switch r.Method {
	case http.MethodGet:
		handleNetworksGet(w, userID)
	case http.MethodPost:
		handleNetworksPost(w, r, userID, username)
	case http.MethodPut:
		fmt.Println("PUT /networks")
	case http.MethodDelete:
		fmt.Println("DELETE /networks")
	default:
		ErrorMethodNotAllowed(w, r)
	}
}

func handleNetworksGet(w http.ResponseWriter, userID uint) {
	fmt.Println("GET /networks")
	var networkUsers []NetworkUser
	DB.Preload("Network.Users").Where(&NetworkUser{UserId: userID}).Find(&networkUsers)
	WriteJsonResponse(w, ConvertNetworkUserList(networkUsers))
}

func handleNetworksPost(w http.ResponseWriter, r *http.Request, userID uint, username string) {
	fmt.Println("POST /networks")
	var networkDTO DTONetwork
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&networkDTO)
	if err != nil {
		ErrorBadRequest(w, r, err)
		return
	}

	network := Network{
		Name:   networkDTO.Name,
		UserId: userID,
	}
	DB.Create(&network)

	addedOwner := false
	for _, uDTO := range networkDTO.Members {
		var user User
		if uDTO.Username == username {
			addedOwner = true
		}
		if DB.Where(&User{Username: uDTO.Username}).First(&user).RecordNotFound() {
			continue
		}

		networkUser := NetworkUser{
			UserId:    user.ID,
			NetworkId: network.ID,
			ColorHex:  networkDTO.ColorHex,
		}
		DB.Create(&networkUser)
	}

	if !addedOwner {
		var user User
		DB.Where(&User{}, userID).Find(&user)
		networkUser := NetworkUser{
			UserId:    user.ID,
			NetworkId: network.ID,
			ColorHex:  networkDTO.ColorHex,
		}
		DB.Create(&networkUser)
	}

	var networkFinal Network
	DB.Preload("Users").Where(&Network{}, network.ID).Find(&networkFinal)
	WriteJsonResponseWithStatus(w, ConvertNetwork(networkFinal, networkDTO.ColorHex), http.StatusCreated)
}

func NotifyAllNetworkMembers(networkID uint, message string) {
	var network Network
	if DB.Preload("Users").First(&network, networkID).RecordNotFound() {
		return
	}

	smsMessage := fmt.Sprintf("[%s] %s", network.Name, message)

	count := 0
	for _, user := range network.Users {
		if user.MobilePhone == "" {
			continue
		}
		_, errConv := strconv.Atoi(user.MobilePhone)
		if errConv != nil || len(user.MobilePhone) != 10 {
			continue
		}

		if err := SendSMS(user.MobilePhone, smsMessage); err != nil {
			count++
		}
	}
}
