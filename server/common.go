package server

import (
	"net/http"
)

func AuthenticateRequest(r *http.Request) (bool, map[interface{}]interface{}) {
	session, _ := Store.Get(r, "session_calendays")
	return !session.IsNew, session.Values
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
