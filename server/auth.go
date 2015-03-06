package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
	"github.com/localhots/empact/job"
	"github.com/localhots/empact/task"
)

func authHelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	helloTmpl.ExecuteTemplate(w, "hello", map[string]interface{}{})
}

func authSigninHandler(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Set("client_id", config.C().ClientID)
	params.Set("redirect_uri", config.C().RedirectURI)
	params.Set("scope", "read:org, repo, admin:org_hook")
	http.Redirect(w, r, config.C().AuthURL+"?"+params.Encode(), 302)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("error") != "" {
		w.Write([]byte(r.FormValue("error_description")))
	} else {
		code := r.FormValue("code")
		fmt.Println("Got code: ", code)

		res := make(chan string)
		job.Enqueue(&task.FetchAccessTokenTask{
			Code:   code,
			Result: res,
			Task:   &db.Task{},
		})

		if login, ok := <-res; ok {
			authorize(r, login)
		} else {
			panic("Failed to access token or user info")
		}
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if currentUser(r) == "" {
		http.Redirect(w, r, "/auth/hello", 302)
	}
}