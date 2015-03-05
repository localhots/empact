package task

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
)

type (
	FetchAccessTokenTask struct {
		Code   string
		Result chan string
		*db.Task
	}
)

func FetchAccessToken(tk Tasker) {
	t := tk.(*FetchAccessTokenTask)
	payload := url.Values{}
	payload.Set("client_id", config.C().ClientID)
	payload.Set("client_secret", config.C().ClientSecret)
	payload.Set("code", t.Code)
	payload.Set("redirect_uri", config.C().RedirectURI)

	buf := bytes.NewBuffer([]byte(payload.Encode()))
	resp, err := http.Post(config.C().AccessTokenURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	pairs, _ := url.ParseQuery(string(body))

	t.Result <- pairs.Get("access_token")
}
