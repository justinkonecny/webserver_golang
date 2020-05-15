package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var store = sessions.NewFilesystemStore("/Users/justinkonecny/Documents/Git Workspace/go_server/fs", []byte("SECRET_KEY"))

func main() {
	store.Options = &sessions.Options{
		MaxAge: 10,
		HttpOnly: true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handleLogin)

	log.Fatal(http.ListenAndServe("localhost:8081", mux))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_calendays")
	if session.IsNew {
		headerUsername := r.Header.Values("username")
		headerPassword := r.Header.Values("password")
		if len(headerUsername) != 1 || len(headerPassword) != 1 {
			http.Error(w, "Missing/incorrect authentication headers", http.StatusBadRequest)
		}

		username := headerUsername[0]
		password := headerPassword[0]

		fmt.Println("Saving user credentials...")
		fmt.Println("Username:", username)
		fmt.Println("Password:", password)

		session.Values["username"] = username
		session.Values["password"] = password
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		username := session.Values["username"]
		password := session.Values["password"]

		fmt.Println("Loading existing session")
		fmt.Println("Username:", username)
		fmt.Println("Password:", password)
	}
}
