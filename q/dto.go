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

type UserLoginRequest struct {
	SpotifyUserID string `json:"spotifyUserID"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}

type UserResponse struct {
	ID            uint   `json:"id"`
	UUID          string `json:"uuid"`
	SpotifyUserID string `json:"spotifyUserID"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}

func GetResponseFromUser(user User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		UUID:          user.UUID,
		SpotifyUserID: user.SpotifyUserID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
	}
}
