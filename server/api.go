package server

import (
	"encoding/json"
	"net/http"

	"github.com/localhots/empact/db"
)

func apiOrgsHandler(w http.ResponseWriter, r *http.Request) {
	login := currentUser(r)
	orgs := db.UserOrgs(login)
	b, _ := json.Marshal(orgs)

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(b)
}

func apiTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams := db.OrgTeams(r.FormValue("org"))
	b, _ := json.Marshal(teams)

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(b)
}

func apiReposHandler(w http.ResponseWriter, r *http.Request) {
	repos := db.OrgRepos(r.FormValue("org"))
	b, _ := json.Marshal(repos)

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(b)
}
