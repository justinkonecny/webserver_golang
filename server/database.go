package server

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"time"
)

type User struct {
	gorm.Model
	FirebaseUuid string
	FirstName    string
	LastName     string
	Email        string
	Networks     []Network `gorm:"many2many:network_user"`
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

type NetworkUser struct {
	gorm.Model
	UserId    uint
	NetworkId uint
	ColorHex  string
	User      User
	Network   Network
}

var DB *gorm.DB

func InitDatabase() {
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
		fmt.Println("Error connecting to DB")
		fmt.Println(err)
		panic("Failed to connect to database")
	}
	DB = database
	DB.SingularTable(true)
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Network{})
	DB.AutoMigrate(&Event{})
	DB.AutoMigrate(&NetworkUser{})
	fmt.Println("Successfully initialized database connection")
}
