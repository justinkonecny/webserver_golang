package q

import (
	"fmt"
	"github.com/gorilla/sessions"
	"os"
	"sync"
)

var Store *sessions.FilesystemStore

const cookieAgeMinutes = 525600 // 1 year

const KeyUserID = "UserID"

func InitStore(wg *sync.WaitGroup) {
	defer wg.Done()
	storeSecretKey := os.Getenv("STORE_Q_SK")
	if storeSecretKey == "" {
		panic("Missing Q session store secret key")
	}

	Store = sessions.NewFilesystemStore("q", []byte(storeSecretKey))
	Store.Options = &sessions.Options{MaxAge: cookieAgeMinutes * 60, HttpOnly: true}
	fmt.Println("Successfully initialized Q session Store")
}
