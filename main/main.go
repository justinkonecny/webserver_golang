package main

import (
	"../server"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const serverPort = 8081

func main() {
	server.InitDatabase()
	defer server.DB.Close()
	server.InitStore()
	InitWebServer()
}

func InitWebServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", server.HandleLogin)
	mux.HandleFunc("/events", server.HandleEvents)
	mux.HandleFunc("/networks", server.HandleNetworks)

	fmt.Printf("Starting web server on port %d...\n", serverPort)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(serverPort), mux))
}
