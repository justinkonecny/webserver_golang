package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"os"
	"sync"
)

var Store *sessions.FilesystemStore

const cookieAgeMinutes = 60

const KeyFirebaseUUID = "FirebaseUUID"
const KeyUserID = "UserID"
const KeyUserEmail = "UserEmail"
const KeyUsername = "Username"

func InitStore(wg *sync.WaitGroup) {
	defer wg.Done()
	storeSecretKey := os.Getenv("STORE_SK")
	if storeSecretKey == "" {
		panic("Missing session Store secret key")
	}

	Store = sessions.NewFilesystemStore("", []byte(storeSecretKey))
	Store.Options = &sessions.Options{MaxAge: cookieAgeMinutes * 60, HttpOnly: true}
	fmt.Println("Successfully initialized session Store")
}
