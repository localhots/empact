package server

import (
	"net/http"

	"github.com/localhots/empact/db"
)

func apiUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := parseRequest(w, r)
	orgs := db.UserOrgs(req.login)
	req.respondWith(orgs)
}

func apiOrgTeamsHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	teams := db.OrgTeams(stat.Org)
	req.respondWith(teams)
}

func apiOrgUsersHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	users := db.OrgUsers(stat.Org)
	req.respondWith(users)
}

func apiOrgReposHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	repos := db.OrgRepos(stat.Org)
	req.respondWith(repos)
}

func apiOrgWeekRangeHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	min, max := db.OrgWeekRange(stat.Org)
	req.respondWith([]int{min, max})
}
