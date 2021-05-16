package q

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"sync"
)

type User struct {
	gorm.Model
	SpotifyUserID string
	Email         string
	FirstName     string
	LastName      string
	PasswordHash  string
}

type UserFriend struct {
	gorm.Model
	UserId       uint
	FriendUserID uint
	User         User
	UserFriend   User `gorm:"foreignKey:FriendUserId"`
}

var DB *gorm.DB

func InitDatabase(wg *sync.WaitGroup) {
	defer wg.Done()
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_Q_USER")
	dbPassword := os.Getenv("DB_Q_PASSWORD")
	dbSchema := os.Getenv("DB_Q_SCHEMA")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbSchema == "" {
		panic("Missing Q database connection information")
	}

	dbInfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbSchema)
	database, err := gorm.Open("mysql", dbInfo)
	if err != nil {
		panic("Failed to connect to Q database")
	}

	DB = database
	DB.SingularTable(true)
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&UserFriend{})
	fmt.Println("Successfully initialized Q database connection")
}
