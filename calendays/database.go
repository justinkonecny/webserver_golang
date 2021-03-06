package calendays

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"sync"
	"time"
)

type User struct {
	gorm.Model
	FirebaseUuid string
	FirstName    string
	LastName     string
	Email        string
	Username     string
	MobilePhone  string
	Networks     []Network `gorm:"many2many:network_user"`
}

type Network struct {
	gorm.Model
	Name   string
	UserId uint
	User   User
	Users  []User `gorm:"many2many:network_user"`
}

type Event struct {
	gorm.Model
	StartDate time.Time
	EndDate   time.Time
	Name      string
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

func InitDatabase(wg *sync.WaitGroup) {
	defer wg.Done()
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_CAL_USER")
	dbPassword := os.Getenv("DB_CAL_PASSWORD")
	dbSchema := os.Getenv("DB_CAL_SCHEMA")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbSchema == "" {
		panic("Missing Calendays database connection information")
	}

	dbInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbSchema)
	database, err := gorm.Open("mysql", dbInfo)
	if err != nil {
		panic("Failed to connect to Calendays database")
	}

	DB = database
	DB.SingularTable(true)
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Network{})
	DB.AutoMigrate(&Event{})
	DB.AutoMigrate(&NetworkUser{})
	fmt.Println("Successfully initialized Calendays database connection")
}
