package server

import (
	"net/http"

	"github.com/localhots/empact/db"
)

func statOrgReposTop(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatOrgReposTop(stat.org, stat.from, stat.to)
	req.respondWith(top)
}

func statOrgReposActivity(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatOrgReposActivity(stat.org, stat.from, stat.to)
	req.respondWith(activity)
}

func statOrgTeamsTop(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatOrgTeamsTop(stat.org, stat.from, stat.to)
	req.respondWith(top)
}

func statOrgTeamsActivity(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatOrgTeamsActivity(stat.org, stat.from, stat.to)
	req.respondWith(activity)
}

func statOrgUsersTop(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatOrgUsersTop(stat.org, stat.from, stat.to)
	req.respondWith(top)
}
