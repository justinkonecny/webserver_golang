package q

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserSignupRequest struct {
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Password      string `json:"password"`
}

type UserLoginRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
}

type UserResponse struct {
	ID            uint   `json:"id"`
	UUID          string `json:"uuid"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
}

func GetResponseFromUser(user User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		UUID:          user.UUID,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
	}
}
