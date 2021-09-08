package q

import (
	"../common"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
		userSignupRequest.Password == "" {
		errMsg := ErrorResponse{Error: "No field may be blank"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusBadRequest)
		return
	}

	// Check for duplicate users
	var users []User
	DB.Where("email = ?", userSignupRequest.Email).Find(&users)
	if len(users) != 0 {
		fmt.Println("(Q-A.HS) User sign up failed!")
		errMsg := ErrorResponse{Error: "User already exists"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusConflict)
		return
	}

	passwordHash, err := common.HashPassword(userSignupRequest.Password)
	if err != nil {
		fmt.Println("(Q-A.HS) Failed to hash password!")
		fmt.Println(err)
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Create new User in DB
	newUser := User{
		Email:         userSignupRequest.Email,
		FirstName:     userSignupRequest.FirstName,
		LastName:      userSignupRequest.LastName,
		PasswordHash:  passwordHash,
		UUID:          uuid.New().String(),
	}
	if err := DB.Create(&newUser).Error; err != nil {
		fmt.Println("(Q-A.HS) User record creation failed with error:", err)
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	if session := CreateNewSession(w, r, newUser); session == nil {
		fmt.Println("(Q-A.HS) User session creation failed!")
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to create a new User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Return a 201 Created
	response := GetResponseFromUser(newUser)
	common.WriteJsonResponseWithStatus(w, response, http.StatusCreated)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	common.EnableCORS(w, r)
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}
	fmt.Println("Q POST /login")

	// Decode the JSON body
	var userLoginRequest UserLoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userLoginRequest)
	if err != nil {
		common.ErrorBadRequest(w, r, err)
		return
	}

	if userLoginRequest.Email == "" ||
		userLoginRequest.Password == "" {
		errMsg := ErrorResponse{Error: "No field may be blank"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusBadRequest)
		return
	}

	var user User
	if result := DB.First(&user, "email = ?", userLoginRequest.Email); result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !common.VerifyPassword(user.PasswordHash, userLoginRequest.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if session := CreateNewSession(w, r, user); session == nil {
		fmt.Println("(Q-A.HL) User session creation failed!")
		errMsg := ErrorResponse{Error: "Unknown error occurred trying to authenticate User"}
		common.WriteJsonResponseWithStatus(w, errMsg, http.StatusInternalServerError)
		return
	}

	response := GetResponseFromUser(user)
	common.WriteJsonResponseWithStatus(w, response, http.StatusOK)
}
