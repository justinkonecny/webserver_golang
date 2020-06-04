package main

import (
	"../server"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting application...")

	server.InitDatabase()
	defer server.DB.Close()
	server.InitStore()
	InitWebServer()
}

func InitWebServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server.SetupCommon()
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/login", server.HandleLogin)
	router.HandleFunc("/signup", server.HandleSignup)

	router.HandleFunc("/status/user", server.HandleStatusUser)

	router.HandleFunc("/events", server.HandleEvents)
	router.HandleFunc("/networks", server.HandleNetworks)
	router.HandleFunc("/users", server.HandleUsers)

	fmt.Printf("Starting web server on port %s...\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	server.EnableCORS(w, r)
}
