package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var isDevEnv bool
//var SNSNeverSubscribedID uint
var SNSSubscribedID uint
//var SNSUnsubscribedID uint

func SetupCommon() {
	dev, err := strconv.ParseBool(os.Getenv("DEV"))
	if err != nil {
		isDevEnv = false
	} else {
		isDevEnv = dev
	}

	//SNSNeverSubscribedID = 1
	SNSSubscribedID = 2
	//SNSUnsubscribedID = 3
}

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (bool, map[interface{}]interface{}) {
	EnableCORS(w, r)
	session, _ := Store.Get(r, "session_calendays")
	return !session.IsNew, session.Values
}

func EnableCORS(w http.ResponseWriter, r *http.Request) {
	origins := [3]string{
		"https://calendays.jkonecny.com",
		"https://www.calendays.jkonecny.com",
		"https://www.jkonecny.com",
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
