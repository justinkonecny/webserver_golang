package calendays

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
	storeSecretKey := os.Getenv("STORE_CAL_SK")
	if storeSecretKey == "" {
		panic("Missing Calendays session store secret key")
	}

	Store = sessions.NewFilesystemStore("calendays", []byte(storeSecretKey))
	Store.Options = &sessions.Options{MaxAge: cookieAgeMinutes * 60, HttpOnly: true}
	fmt.Println("Successfully initialized Calendays session Store")
}
