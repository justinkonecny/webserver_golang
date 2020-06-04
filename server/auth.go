package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

type StatusUser struct {
	Username string
	Email    string
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w, r)
		return
	}

	fmt.Println("POST /login")

	var session *sessions.Session
	headerFirebaseUUID := r.Header.Values(KeyFirebaseUUID)
	if len(headerFirebaseUUID) == 1 {
		fmt.Println("New")
		session, _ = Store.Get(r, "session_calendays")
		if !processNewSession(session, w, r) {
			fmt.Println("(AE01): Error creating new session")
			return // Something went wrong creating a new session
		}
	} else {
		fmt.Println("Existing")
		session, _ = Store.Get(r, "session_calendays")
		if session.IsNew && !processNewSession(session, w, r) {
			fmt.Println("(AE02): Error creating new session")
			return // Something went wrong creating a new session
		}
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

		fmt.Println("(AE03) User sign up failed!")
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

	if err := DB.Create(&newUser).Error; err != nil {
		fmt.Println("(AE04) User sign record creation failed!")
		fmt.Println(err)
		statusError := map[string]bool{
			"ExistingEmail":        false,
			"ExistingUsername":     false,
			"ExistingFirebaseUUID": false,
			"Success":              false,
		}
		WriteJsonResponseWithStatus(w, statusError, http.StatusInternalServerError)
		return
	}

	statusSuccess := map[string]bool{
		"ExistingEmail":        false,
		"ExistingUsername":     false,
		"ExistingFirebaseUUID": false,
		"Success":              true,
	}

	WriteJsonResponseWithStatus(w, statusSuccess, http.StatusCreated)
}

func HandleStatusUser(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w, r)
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w, r)
		return
	}

	fmt.Println("POST /status/user")
	var statusUser StatusUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&statusUser)
	if err != nil {
		ErrorBadRequest(w, r, err)
		return
	}

	usernameNotFound := DB.Where(&User{Username: statusUser.Username}).First(&User{}).RecordNotFound()
	emailNotFound := DB.Where(&User{Email: statusUser.Email}).First(&User{}).RecordNotFound()

	status := map[string]bool{
		"ExistingUsername": !usernameNotFound,
		"ExistingEmail":    !emailNotFound,
	}
	WriteJsonResponse(w, status)
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
	session.Values[KeyUsername] = user.Username
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	fmt.Printf("Successfully authenticated user '%s'\n", user.FirstName)
	return true
}
