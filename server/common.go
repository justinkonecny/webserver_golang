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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func ErrorMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func ErrorUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func ErrorInternalServerError(w http.ResponseWriter) {
	http.Error(w, "Oops, something went wrong", http.StatusInternalServerError)
}
