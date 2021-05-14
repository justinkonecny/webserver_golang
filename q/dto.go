package q

type ErrorMessage struct {
	Error string `json:"error"`
}

type DTOUserSignup struct {
	SpotifyUserID string `json:"spotifyUserID"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}
