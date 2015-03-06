package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
)

func FetchAccessToken(code string, result chan string) {
	payload := url.Values{}
	payload.Set("client_id", config.C().ClientID)
	payload.Set("client_secret", config.C().ClientSecret)
	payload.Set("code", code)
	payload.Set("redirect_uri", config.C().RedirectURI)

	buf := bytes.NewBuffer([]byte(payload.Encode()))
	fmt.Println("Requesting token")
	resp, err := http.Post(config.C().AccessTokenURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	pairs, _ := url.ParseQuery(string(body))
	token := pairs.Get("access_token")
	fmt.Println("Got token: " + token)

	fmt.Println("Requesting info")
	resp, err = http.Get("https://api.github.com/user?access_token=" + token)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))

	user := &db.User{}
	json.Unmarshal(body, &user)
	user.Save()
	fmt.Println("Saving user", user)

	tok := &db.Token{
		User:  user.Login,
		Token: token,
	}
	fmt.Println("Saving token", tok)
	tok.Save()

	result <- user.Login
}
