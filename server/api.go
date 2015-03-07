package server

import (
	"net/http"

	"github.com/localhots/empact/db"
)

func apiOrgsHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := parseRequest(w, r)
	orgs := db.UserOrgs(req.login)
	req.respondWith(orgs)
}

func apiTeamsHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	teams := db.OrgTeams(stat.org)
	req.respondWith(teams)
}

func apiReposHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	repos := db.OrgRepos(stat.org)
	req.respondWith(repos)
}
