package server

import (
	"net/http"

	"github.com/localhots/empact/db"
)

func apiOrgsHandler(w http.ResponseWriter, r *http.Request) {
	login := sessionUser(r)
	orgs := db.UserOrgs(login)
	respondWith(w, orgs)
}

func apiTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams := db.OrgTeams(r.FormValue("org"))
	respondWith(w, teams)
}

func apiReposHandler(w http.ResponseWriter, r *http.Request) {
	repos := db.OrgRepos(r.FormValue("org"))
	respondWith(w, repos)
}
