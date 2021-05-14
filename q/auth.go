package q

import (
	"../common"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	common.EnableCORS(w, r)
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}
	fmt.Println("Q POST /signup")

	// Decode the JSON body
	var userDTO DTOUserSignup
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userDTO)
	if err != nil {
		common.ErrorBadRequest(w, r, err)
		return
	}

	// Check for duplicate users
	var users []User
	DB.Where("spotify_user_id = ?", userDTO.SpotifyUserID).Find(&users)
	if len(users) != 0 {
		fmt.Println("(Q-AHS) User sign up failed!")
		errMsg := ErrorMessage{Error: "User already exists"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusConflict)
		return
	}

	// Create new User in DB
	newUser := User{
		SpotifyUserID: userDTO.SpotifyUserID,
		FirstName:     userDTO.FirstName,
		LastName:      userDTO.LastName,
	}
	if err := DB.Create(&newUser).Error; err != nil {
		fmt.Println("(Q-AHS) User record creation failed!")
		fmt.Println(err)
		errMsg := ErrorMessage{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Return a 201 Created
	w.WriteHeader(http.StatusCreated)
}
