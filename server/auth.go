package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorMethodNotAllowed(w)
		return
	}

	session, _ := Store.Get(r, "session_calendays")
	if session.IsNew && !processNewSession(session, w, r) {
		return // Something went wrong creating a new session
	}

	firebaseUUID := session.Values[KeyFirebaseUUID]
	userID := session.Values[KeyUserID]
	fmt.Println("UUID:", firebaseUUID)
	fmt.Println("UserID:", userID)
	EnableCORS(w)
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
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	fmt.Printf("Successfully authenticated user '%s'\n", user.FirstName)
	return true
}
