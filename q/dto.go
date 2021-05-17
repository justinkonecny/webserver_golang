package q

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserSignupRequest struct {
	SpotifyUserID string `json:"spotifyUserID"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Password      string `json:"password"`
}

type UserResponse struct {
	Id            uint   `json:"id"`
	SpotifyUserID string `json:"spotifyUserID"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}
