package main

import (
	"../calendays"
	"../common"
	"../libertycars"
	"../q"
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

	common.SetupCommon()

	var wg sync.WaitGroup
	wg.Add(6)

	// initialize Calendays
	go calendays.InitStore(&wg)
	go calendays.InitDatabase(&wg)

	// initialize Q
	go q.InitAuthDetails(&wg)
	go q.InitStore(&wg)
	go q.InitDatabase(&wg)

	// initialize LibertyCars
	go libertycars.InitAuth(&wg)

	wg.Wait()

	defer calendays.DB.Close()
	defer q.DB.Close()

	StartWebServer()
}

func StartWebServer() {
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
	router.HandleFunc("/", handleHome)

	// define Calendays routes
	routerCalendays := router.PathPrefix("/c").Subrouter()
	routerCalendays.HandleFunc("/login", calendays.HandleLogin)
	routerCalendays.HandleFunc("/signup", calendays.HandleSignup)
	routerCalendays.HandleFunc("/events", calendays.HandleEvents)
	routerCalendays.HandleFunc("/networks", calendays.HandleNetworks)
	routerCalendays.HandleFunc("/users", calendays.HandleUsers)
	routerCalendays.HandleFunc("/status/user", calendays.HandleStatusUser)

	// define Q routes
	routerQ := router.PathPrefix("/q").Subrouter()
	q.DefineRoutes(routerQ)

	// define Liberty Cars routes
	routerLibertyCars := router.PathPrefix("/lc").Subrouter()
	routerLibertyCars.HandleFunc("/search", libertycars.HandleSearch)
	routerLibertyCars.HandleFunc("/listing", libertycars.HandleListing)

	fmt.Printf("Starting web server on port %s...\n", port)

	if isDevEnv {
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
	} else {
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+port, "certs/fullchain.pem", "certs/privkey.pem", router))
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	common.EnableCORS(w, r)
}
