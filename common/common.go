package common

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password")
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func WriteJsonResponse(w http.ResponseWriter, data interface{}) bool {
	return WriteJsonResponseWithStatus(w, data, http.StatusOK)
}

func WriteJsonResponseWithStatus(w http.ResponseWriter, data interface{}, status int) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		ErrorInternalServerError(w)
		return false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, jsonErr := w.Write(jsonData)
	return jsonErr == nil
}

func EnableCORS(w http.ResponseWriter, r *http.Request) {
	allowedOrigin := "https://jkonecny.com"
	if isDevEnv {
		devOrigins := [2]string{
			"http://localhost:3000",
			"http://localhost:8080",
		}

		for _, origin := range devOrigins {
			if r.Header.Get("Origin") == origin {
				allowedOrigin = origin
				fmt.Printf("Matching origin: '%s'\n", origin)
				break
			}
		}

	} else {
		origins := [6]string{
			"https://calendays.jkonecny.com",
			"https://www.calendays.jkonecny.com",
			"https://libertycars.jkonecny.com",
			"https://www.libertycars.jkonecny.com",
			"https://www.jkonecny.com",
		}

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
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, FirebaseUUID, LC-API-Key")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
