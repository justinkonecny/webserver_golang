package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var isDevEnv bool

func SetupCommon() {
	dev, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		isDevEnv = false
	} else {
		isDevEnv = dev
	}
}

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (bool, map[interface{}]interface{}) {
	EnableCORS(w, r)
	session, _ := Store.Get(r, "session_calendays")
	return !session.IsNew, session.Values
}

func EnableCORS(w http.ResponseWriter, r *http.Request) {
	origins := [2]string{
		"https://calendays.jkonecny.com",
		"https://www.calendays.jkonecny.com",
	}

	allowedOrigin := "https://jkonecny.com"
	if isDevEnv {
		allowedOrigin = "http://localhost:3000"
	} else {
		for _, origin := range origins {
			if r.Header.Get("Origin") == origin {
				allowedOrigin = origin
				fmt.Printf("Matching origin: '%s'\n", origin)
				break
			}
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
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

func ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("Unauthorized"))
}

func ErrorInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("Oops, something went wrong"))
}

func ErrorBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("Invalid request body: " + err.Error()))
}
