package q

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"os"
	"sync"
)

var clientID string
var clientSecret string
var authHeader string
var clientCallback = "q://spotify-login-callback"
var spotifyAccountsEndpoint = "https://accounts.spotify.com/api/token"

func InitAuthDetails(wg *sync.WaitGroup) {
	defer wg.Done()

	clientID = os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		panic("Missing Spotify client ID")
	}
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		panic("Missing Spotify client secret")
	}

	authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret))

	fmt.Println("Successfully initialized Q authentication")
}

func DefineRoutes(router *mux.Router) {
	router.HandleFunc("/spotify/token", HandleSpotifyAuthToken)
	router.HandleFunc("/spotify/refresh_token", HandleSpotifyRefreshToken)

	router.HandleFunc("/auth/signup", HandleSignup)
	router.HandleFunc("/auth/login", HandleLogin)
}
