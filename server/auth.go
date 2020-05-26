package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

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

	if len(users) != 1 {
		_, _ = w.Write([]byte("User already exists"))
		w.WriteHeader(http.StatusConflict)
		return
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