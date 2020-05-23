package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
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
	fmt.Println("UUID:", firebaseUUID)
	fmt.Println("UserID:", userID)
	fmt.Println("UserEmail:", userEmail)
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
