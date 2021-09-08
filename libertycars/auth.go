package libertycars

import (
	"../common"
	"net/http"
	"os"
)

var LC_API_KEY string

func InitAuth() {
	LC_API_KEY = os.Getenv("LC_API_KEY")

	if LC_API_KEY == "" {
		panic("Missing LibertyCars API Key")
	}
}

func AuthenticateLibertyCarsRequest(w http.ResponseWriter, r *http.Request) bool {
	common.EnableCORS(w, r)
	apiKey := r.Header.Get("LC_API_KEY")

	if apiKey == "" || apiKey != LC_API_KEY {
		common.ErrorUnauthorized(w, r)
		return false
	}

	return true
}
