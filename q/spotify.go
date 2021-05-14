package q

import (
	"../common"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func HandleSpotifyAuthToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	values, err := url.ParseQuery(string(body))
	if err != nil || values == nil {
		log.Println("Something went wrong processing 'code' request")
		common.ErrorBadRequest(w, r, fmt.Errorf("missing 'code' parameter"))
		return
	}

	code := values.Get("code")
	if code == "" {
		log.Println("Missing 'code' parameter")
		common.ErrorBadRequest(w, r, fmt.Errorf("missing 'code' parameter"))
		return
	}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", clientCallback)
	form.Add("code", code)

	req, _ := http.NewRequest("POST", spotifyAccountsEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authHeader)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Token request error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Token (%s): %s\n", resp.Status, respBody)

	if resp.StatusCode != http.StatusOK {
		common.ErrorBadRequest(w, r, fmt.Errorf("something went wrong"))
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, _ = w.Write(respBody)
}

func HandleSpotifyRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.ErrorMethodNotAllowed(w, r)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	values, err := url.ParseQuery(string(body))
	if err != nil || values == nil {
		log.Println("Something went wrong processing 'refresh_token' request")
		common.ErrorBadRequest(w, r, fmt.Errorf("missing 'refresh_token' parameter"))
		return
	}

	refreshToken := values.Get("refresh_token")
	if refreshToken == "" {
		log.Println("Missing 'refresh_token' parameter")
		common.ErrorBadRequest(w, r, fmt.Errorf("missing 'refresh_token' parameter"))
		return
	}

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", spotifyAccountsEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authHeader)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Refresh request error: %s\n", err.Error())
		common.ErrorBadRequest(w, r, err)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Refresh (%s): %s\n", resp.Status, respBody)

	if resp.StatusCode != http.StatusOK {
		common.ErrorBadRequest(w, r, fmt.Errorf("something went wrong"))
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	_, _ = w.Write(respBody)
}

func printRequest(r *http.Request) {
	fmt.Printf("\n%v %v %v\n", r.Method, r.URL, r.Proto)

	fmt.Printf("Host: %v\n", r.Host)

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Printf("%v: %v\n", name, h)
		}
	}
	fmt.Println()
	for name, headers := range r.URL.Query() {
		for _, h := range headers {
			fmt.Printf("%v: %v\n", name, h)
		}
	}
}
