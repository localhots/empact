package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
)

func Authenticate(code string) (token, login string, err error) {
	defer report(time.Now(), "Authenticate")
	if token, err = FetchAccessToken(code); err != nil {
		return
	}
	log.Printf("Got token %q for code %q\n", token, code)

	var user *db.User
	if user, err = FetchUserInfoWithToken(token); err != nil {
		return
	}
	login = user.Login
	db.Later(func() { user.Save() })

	tok := &db.Token{
		User:  login,
		Token: token,
	}
	db.Later(func() { tok.Save() })

	return
}

func FetchAccessToken(code string) (token string, err error) {
	defer report(time.Now(), "FetchAccessToken")
	payload := url.Values{}
	payload.Set("client_id", config.C().ClientID)
	payload.Set("client_secret", config.C().ClientSecret)
	payload.Set("code", code)
	payload.Set("redirect_uri", config.C().RedirectURI)

	log.Printf("Requesting token for code %q", code)
	buf := bytes.NewBuffer([]byte(payload.Encode()))
	var resp *http.Response
	if resp, err = http.Post(config.C().AccessTokenURL, "application/x-www-form-urlencoded", buf); err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	pairs, _ := url.ParseQuery(string(body))
	if token = pairs.Get("access_token"); token == "" {
		err = fmt.Errorf("Failed to fetch access token usign code %q: %s", code, pairs.Get("error_description"))
	}

	return
}

func FetchUserInfoWithToken(token string) (u *db.User, err error) {
	defer report(time.Now(), "FetchUserInfoWithToken")
	var resp *http.Response
	if resp, err = http.Get("https://api.github.com/user?access_token=" + token); err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	u = &db.User{}
	err = json.Unmarshal(body, &u)

	return
}
