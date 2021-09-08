package q

import (
	"../common"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"sync"
)

var Store *sessions.FilesystemStore

const cookieAgeMinutes = 525600 // 1 year

func InitStore(wg *sync.WaitGroup) {
	defer wg.Done()
	storeSecretKey := os.Getenv("STORE_Q_SK")
	if storeSecretKey == "" {
		panic("Missing Q session store secret key")
	}

	Store = sessions.NewFilesystemStore("sessions/q", []byte(storeSecretKey))
	Store.Options = &sessions.Options{MaxAge: cookieAgeMinutes * 60, HttpOnly: true}
	fmt.Println("Successfully initialized Q session Store")
}

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (bool, UserResponse) {
	common.EnableCORS(w, r)
	session, _ := Store.Get(r, "session_q")
	if session.IsNew {
		return false, UserResponse{}
	}

	userResponse := UserResponse{
		ID:            session.Values["ID"].(uint),
		UUID:          session.Values["UUID"].(string),
		Email:         session.Values["Email"].(string),
		FirstName:     session.Values["FirstName"].(string),
		LastName:      session.Values["LastName"].(string),
	}

	return !session.IsNew, userResponse
}

func CreateNewSession(w http.ResponseWriter, r *http.Request, user User) *sessions.Session {
	fmt.Println("(Q-S.CNS) Creating new user session")

	session, err := Store.Get(r, "session_q")
	if err != nil {
		fmt.Println("(Q-S.CNS) Error creating new session:", err)
		return nil
	}
	if !session.IsNew {
		fmt.Println("(Q-S.CNS) Session already exists")
		return nil
	}

	session.Values["ID"] = user.ID
	session.Values["UUID"] = user.UUID
	session.Values["Email"] = user.Email
	session.Values["FirstName"] = user.FirstName
	session.Values["LastName"] = user.LastName

	if err := session.Save(r, w); err != nil {
		fmt.Println("(Q-S.CNS) Failed to save session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	fmt.Println("(Q-S.CNS) Successfully created new session for user")
	return session
}
