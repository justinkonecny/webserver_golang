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
	var userSignupRequest UserSignupRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userSignupRequest)
	if err != nil {
		common.ErrorBadRequest(w, r, err)
		return
	}

	if userSignupRequest.FirstName == "" ||
		userSignupRequest.LastName == "" ||
		userSignupRequest.Email == "" ||
		userSignupRequest.SpotifyUserID == "" ||
		userSignupRequest.Password == "" {
		errMsg := ErrorResponse{Error: "No field may be blank"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusBadRequest)
		return
	}

	// Check for duplicate users
	var users []User
	DB.Where("spotify_user_id = ?", userSignupRequest.SpotifyUserID).Find(&users)
	if len(users) != 0 {
		fmt.Println("(Q-AHS) User sign up failed!")
		errMsg := ErrorResponse{Error: "User already exists"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusConflict)
		return
	}

	passwordHash, err := common.HashPassword(userSignupRequest.Password)
	if err != nil {
		fmt.Println("(Q-AHS) Failed to hash password!")
		fmt.Println(err)
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Create new User in DB
	newUser := User{
		SpotifyUserID: userSignupRequest.SpotifyUserID,
		Email:         userSignupRequest.Email,
		FirstName:     userSignupRequest.FirstName,
		LastName:      userSignupRequest.LastName,
		PasswordHash:  passwordHash,
	}
	if err := DB.Create(&newUser).Error; err != nil {
		fmt.Println("(Q-AHS) User record creation failed!")
		fmt.Println(err)
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Return a 201 Created
	common.WriteJsonResponseWithStatus(w, UserResponse{
		Id:            newUser.ID,
		SpotifyUserID: newUser.SpotifyUserID,
		Email:         newUser.Email,
		FirstName:     newUser.FirstName,
		LastName:      newUser.LastName,
	}, http.StatusCreated)
}
