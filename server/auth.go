package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/localhots/empact/config"
)

const (
	authURL        = "https://github.com/login/oauth/authorize"
	accessTokenURL = "https://github.com/login/oauth/access_token"
)

func authSigninHandler(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Set("client_id", config.C().ClientID)
	params.Set("redirect_uri", config.C().RedirectURI)
	params.Set("scope", "repo")
	http.Redirect(w, r, authURL+"?"+params.Encode(), 302)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("error") != "" {
		w.Write([]byte(r.FormValue("error_description")))
	} else {
		fmt.Println("Got code: ", r.FormValue("code"))
		token := getAccessToken(r.FormValue("code"))
		fmt.Println("Got access token: ", token)
		w.Write([]byte(token))
	}
}

func getAccessToken(code string) string {
	payload := url.Values{}
	payload.Set("client_id", config.C().ClientID)
	payload.Set("client_secret", config.C().ClientSecret)
	payload.Set("code", code)
	payload.Set("redirect_uri", config.C().RedirectURI)

	buf := bytes.NewBuffer([]byte(payload.Encode()))
	resp, err := http.Post(accessTokenURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	pairs, _ := url.ParseQuery(string(body))
	return pairs.Get("access_token")
}
