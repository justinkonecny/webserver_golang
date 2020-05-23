package server

import (
	"net/http"
)

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (bool, map[interface{}]interface{}) {
	EnableCORS(w)
	session, _ := Store.Get(r, "session_calendays")
	return !session.IsNew, session.Values
}

func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, FirebaseUUID")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func ErrorMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("Method not allowed"))
}

func ErrorUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func ErrorInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Oops, something went wrong", http.StatusInternalServerError)
}

func ErrorBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, "Invalid request body: " + err.Error(), http.StatusBadRequest)
}
