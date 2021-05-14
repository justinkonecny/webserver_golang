package common

import "net/http"

func ErrorMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("Method not allowed"))
}

func ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("Unauthorized"))
}

func ErrorInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("Oops, something went wrong"))
}

func ErrorBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("Invalid request body: " + err.Error()))
}

func ErrorBadGatewayAWS(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
	_, _ = w.Write([]byte("Error interacting with AWS"))
}