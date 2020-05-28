package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

type StatusUsername struct {
	Username string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w, r)
		return
	}

	session, _ := Store.Get(r, "session_calendays")
	if session.IsNew && !processNewSession(session, w, r) {
		return // Something went wrong creating a new session
	}

	firebaseUUID := session.Values[KeyFirebaseUUID]
	userID := session.Values[KeyUserID]
	userEmail := session.Values[KeyUserEmail]
	username := session.Values[KeyUsername]
	fmt.Println("UUID:", firebaseUUID)
	fmt.Println("UserID:", userID)
	fmt.Println("UserEmail:", userEmail)
	fmt.Println("Username:", username)
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w, r)
		return
	}

	fmt.Println("POST /signup")
	var userDTO DTOUserSignup
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userDTO)
	if err != nil {
		ErrorBadRequest(w, r, err)
		return
	}

	var users []User
	DB.Where("email = ? OR username = ? OR firebase_uuid = ?", userDTO.Email, userDTO.Username, userDTO.FirebaseUUID).Find(&users)

	if len(users) != 0 {
		existingUsername := false
		existingEmail := false
		existingFirebaseUUID := false
		user := users[0]

		if user.Email == userDTO.Email {
			existingEmail = true
		}
		if user.Username == userDTO.Username {
			existingUsername = true
		}
		if user.FirebaseUuid == userDTO.FirebaseUUID {
			existingFirebaseUUID = true
		}

		statusFailed := map[string]bool{
			"ExistingEmail":        existingEmail,
			"ExistingUsername":     existingUsername,
			"ExistingFirebaseUUID": existingFirebaseUUID,
			"Success":              false,
		}

		WriteJsonResponseWithStatus(w, statusFailed, http.StatusConflict)
		return
	}

	newUser := User{
		FirebaseUuid: userDTO.FirebaseUUID,
		FirstName:    userDTO.FirstName,
		LastName:     userDTO.LastName,
		Email:        userDTO.Email,
		Username:     userDTO.Username,
	}

	DB.NewRecord(&newUser)
	statusSuccess := map[string]bool{
		"ExistingEmail":        false,
		"ExistingUsername":     false,
		"ExistingFirebaseUUID": false,
		"Success":              true,
	}

	WriteJsonResponseWithStatus(w, statusSuccess, http.StatusCreated)
}

func HandleUsername(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w, r)
		return
	}

	fmt.Println("POST /status/username")
	var statusUsername StatusUsername
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&statusUsername)
	if err != nil {
		ErrorBadRequest(w, r, err)
		return
	}

	var user User
	if notFound := DB.Where(&User{Username: statusUsername.Username}).First(&user).RecordNotFound(); notFound {
		statusUnique := map[string]bool{
			"ExistingUsername": false,
		}
		WriteJsonResponse(w, statusUnique)
	} else {
		statusConflict := map[string]bool{
			"ExistingUsername": true,
		}
		WriteJsonResponse(w, statusConflict)
	}
}

func processNewSession(session *sessions.Session, w http.ResponseWriter, r *http.Request) bool {
	headerFirebaseUUID := r.Header.Values(KeyFirebaseUUID)
	if len(headerFirebaseUUID) != 1 || len(headerFirebaseUUID[0]) < 10 {
		http.Error(w, "Missing/incorrect authentication headers", http.StatusBadRequest)
		return false
	}

	firebaseUUID := headerFirebaseUUID[0]
	fmt.Println("Creating new user session...")

	var user User
	result := DB.Where(&User{FirebaseUuid: firebaseUUID}).First(&user)
	if result.RecordNotFound() {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return false
	}

	session.Values[KeyFirebaseUUID] = firebaseUUID
	session.Values[KeyUserID] = user.ID
	session.Values[KeyUserEmail] = user.Email
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	fmt.Printf("Successfully authenticated user '%s'\n", user.FirstName)
	return true
}
