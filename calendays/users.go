package calendays

import (
	"../common"
	"fmt"
	"net/http"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateCalendaysRequest(w, r)
	if !auth {
		common.ErrorUnauthorized(w, r)
		return
	}

	userID := values[KeyUserID].(uint)

	switch r.Method {
	case http.MethodGet:
		handleUsersGet(w, userID)
	case http.MethodPost:
		fmt.Println("POST /users")
	case http.MethodPut:
		fmt.Println("PUT /users")
	case http.MethodDelete:
		fmt.Println("DELETE /users")
	default:
		common.ErrorMethodNotAllowed(w, r)
	}
}

func handleUsersGet(w http.ResponseWriter, userID uint) {
	fmt.Println("GET /users")
	var user User
	DB.First(&user, userID)
	common.WriteJsonResponse(w, ConvertUser(user))
}
