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
	if r.FormValue("error") != "" {
		w.Write([]byte(r.FormValue("error_description")))
		return
	}

	code := r.FormValue("code")
	log.Println("Got code: ", code)
	if _, login, err := task.Authenticate(code); err == nil {
		createSession(r, login)
	} else {
		panic(err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if sessionUser(r) == "" {
		http.Redirect(w, r, "/auth/hello", 302)
	}
}
