package main

import (
	"../server"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	fmt.Println("Starting application...")

	var wg sync.WaitGroup
	wg.Add(3)

	go server.InitAWS(&wg)
	go server.InitStore(&wg)
	go server.InitDatabase(&wg)

	wg.Wait()

	defer server.DB.Close()
	InitWebServer()
}

func InitWebServer() {
	isDevEnv := false
	dev, err := strconv.ParseBool(os.Getenv("DEV"))
	if err == nil {
		isDevEnv = dev
	}

	port := os.Getenv("PORT")
	if port == "" {
		if isDevEnv {
			port = "8080"
		} else {
			port = "8443"
		}
	}

	server.SetupCommon()
	router := mux.NewRouter()
	router.HandleFunc("/login", server.HandleLogin)
	router.HandleFunc("/signup", server.HandleSignup)

	router.HandleFunc("/events", server.HandleEvents)
	router.HandleFunc("/networks", server.HandleNetworks)
	router.HandleFunc("/users", server.HandleUsers)

	router.HandleFunc("/status/user", server.HandleStatusUser)

	if !isDevEnv {
		// Don't register these routes in a development environment
		router.HandleFunc("/notifications/subscribe", server.HandleSubscribe)
	}

	router.HandleFunc("/", handleHome)

	fmt.Printf("Starting web server on port %s...\n", port)

	if isDevEnv {
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
	} else {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+port, "fullchain.pem", "privkey.pem", router))

	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	server.EnableCORS(w, r)
}
