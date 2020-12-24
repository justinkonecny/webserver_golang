package main

import (
	"../ios"
	"../libertycars"
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

	server.SetupCommon()

	var wg sync.WaitGroup
	wg.Add(3)

	go ios.InitIOS(&wg)
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

	router := mux.NewRouter()
	router.HandleFunc("/c/login", server.HandleLogin)
	router.HandleFunc("/c/signup", server.HandleSignup)

	router.HandleFunc("/c/events", server.HandleEvents)
	router.HandleFunc("/c/networks", server.HandleNetworks)
	router.HandleFunc("/c/users", server.HandleUsers)

	router.HandleFunc("/c/status/user", server.HandleStatusUser)

	router.HandleFunc("/api/token", ios.HandleToken)
	router.HandleFunc("/api/refresh_token", ios.HandleRefresh)

	router.HandleFunc("/lc/search", libertycars.HandleSearch)
	router.HandleFunc("/lc/listing", libertycars.HandleListing)

	router.HandleFunc("/", handleHome)

	fmt.Printf("Starting web server on port %s...\n", port)

	if isDevEnv {
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
	} else {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+port, "letsencrypt/live/api.jkonecny.com/fullchain.pem", "letsencrypt/live/api.jkonecny.com/privkey.pem", router))
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	server.EnableCORS(w, r)
}
