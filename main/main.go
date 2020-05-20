package main

import (
	"../server"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const serverPort = 8081

func main() {
	fmt.Println("STARTING APP")
	fmt.Printf("%s", os.Getenv("PORT"))

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

	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/login", server.HandleLogin)
	mux.HandleFunc("/events", server.HandleEvents)
	mux.HandleFunc("/networks", server.HandleNetworks)

	fmt.Printf("Starting web server on port %s...\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, mux))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	return
}
