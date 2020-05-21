package main

import (
	"../server"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

const serverPort = 8081

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
		port = strconv.Itoa(serverPort)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", handleHome)
	router.HandleFunc("/login", server.HandleLogin)
	router.HandleFunc("/events", server.HandleEvents)
	router.HandleFunc("/networks", server.HandleNetworks)

	fmt.Printf("Starting web server on port %s...\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
