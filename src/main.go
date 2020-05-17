package main

import (
	"./common"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const serverPort = 8081

func main() {
	common.InitDatabase()
	defer common.DB.Close()
	common.InitStore()
	initWebServer()
}

func initWebServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", common.HandleLogin)
	mux.HandleFunc("/events", common.HandleEvents)

	fmt.Printf("Starting web server on port %d...\n", serverPort)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(serverPort), mux))
}
