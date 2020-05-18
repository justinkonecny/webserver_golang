package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"os"
)

var Store *sessions.FilesystemStore

const cookieAgeMinutes = 60

const KeyFirebaseUUID = "FirebaseUUID"
const KeyUserID = "UserID"

func InitStore() {
	storeSecretKey := os.Getenv("STORE_SK")
	if storeSecretKey == "" {
		panic("Missing session Store secret key")
	}

	Store = sessions.NewFilesystemStore("", []byte(storeSecretKey))
	Store.Options = &sessions.Options{MaxAge: cookieAgeMinutes * 60, HttpOnly: true}
	fmt.Println("Successfully initialized session Store")
}
