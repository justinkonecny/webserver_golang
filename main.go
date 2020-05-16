package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	FirebaseUuid string
	FirstName    string
	LastName     string
	Email        string
}

type Network struct {
	gorm.Model
	Name   string
	UserId uint
	User   User
}

type Event struct {
	gorm.Model
	StartDate time.Time
	EndDate   time.Time
	Location  string
	Message   string
	NetworkId uint
	Network   Network
}

var store *sessions.FilesystemStore
var db *gorm.DB

const SERVER_PORT = 8081

func main() {
	initDatabase()
	defer db.Close()
	initWebServer()
}

func initDatabase() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbSchema := os.Getenv("DB_SCHEMA")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbSchema == "" {
		panic("Missing database connection information")
	}

	dbInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbSchema)
	database, err := gorm.Open("mysql", dbInfo)
	if err != nil {
		panic("Failed to connect to database")
	}

	database.SingularTable(true)
	database.AutoMigrate(&User{})
	database.AutoMigrate(&Network{})
	database.AutoMigrate(&Event{})
	db = database
	fmt.Println("Successfully initialized database connection")
}

func initWebServer() {
	storeSecretKey := os.Getenv("STORE_SK")
	if storeSecretKey == "" {
		panic("Missing session store secret key")
	}

	store = sessions.NewFilesystemStore("/Users/justinkonecny/Documents/Git Workspace/go_server/fs", []byte(storeSecretKey))
	store.Options = &sessions.Options{MaxAge: 10, HttpOnly: true}
	fmt.Println("Successfully initialized session store")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handleLogin)
	fmt.Printf("Starting web server on port %d...", SERVER_PORT)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(SERVER_PORT), mux))
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
