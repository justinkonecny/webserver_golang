package calendays

import (
	"../common"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleNetworks(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateCalendaysRequest(w, r)
	if !auth {
		common.ErrorUnauthorized(w, r)
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
		common.ErrorMethodNotAllowed(w, r)
	}
}

func handleNetworksGet(w http.ResponseWriter, userID uint) {
	fmt.Println("GET /networks")
	var networkUsers []NetworkUser
	DB.Preload("Network.Users").Where(&NetworkUser{UserId: userID}).Find(&networkUsers)
	common.WriteJsonResponse(w, ConvertNetworkUserList(networkUsers))
}

func handleNetworksPost(w http.ResponseWriter, r *http.Request, userID uint, username string) {
	fmt.Println("POST /networks")
	var networkDTO DTONetwork
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&networkDTO)
	if err != nil {
		common.ErrorBadRequest(w, r, err)
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
	common.WriteJsonResponseWithStatus(w, ConvertNetwork(networkFinal, networkDTO.ColorHex), http.StatusCreated)
}
