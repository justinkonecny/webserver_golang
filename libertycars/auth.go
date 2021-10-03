package libertycars

import (
	"../common"
	"net/http"
	"os"
	"sync"
)

var LC_API_KEY string

func InitAuth(wg *sync.WaitGroup) {
	defer wg.Done()

	LC_API_KEY = os.Getenv("LC_API_KEY")

	if LC_API_KEY == "" {
		panic("Missing LibertyCars API Key")
	}
}

func AuthenticateLibertyCarsRequest(w http.ResponseWriter, r *http.Request) bool {
	common.EnableCORS(w, r)
	apiKey := r.Header.Get("LC-API-Key")

	if apiKey == "" || apiKey != LC_API_KEY {
		common.ErrorUnauthorized(w, r)
		return false
	}

	return true
}
