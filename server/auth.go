package server

import (
	"log"
	"net/http"
	"net/url"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/task"
)

func authSigninHandler(w http.ResponseWriter, r *http.Request) {
	params := url.Values{}
	params.Set("client_id", config.C().ClientID)
	params.Set("redirect_uri", config.C().RedirectURI)
	params.Set("scope", "read:org, repo, admin:org_hook")
	http.Redirect(w, r, config.C().AuthURL+"?"+params.Encode(), 302)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := parseRequest(w, r)
	if r.FormValue("error") != "" {
		w.Write([]byte(r.FormValue("error_description")))
		return
	}

	code := r.FormValue("code")
	log.Printf("Got code %q\n", code)

	if token, login, err := task.Authenticate(code); err == nil {
		req.authorize(token, login)
	} else {
		panic(err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if req, _ := parseRequest(w, r); req.login == "" {
		http.Redirect(w, r, "/", 302)
	}
}
