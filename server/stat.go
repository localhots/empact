package server

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/localhots/empact/db"
)

func statOrgTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatOrgTop(req.login, structs.Map(stat))
	req.respondWith(top)
}

func statOrgActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatOrgActivity(req.login, structs.Map(stat))
	req.respondWith(activity)
}

func statTeamTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatTeamTop(req.login, structs.Map(stat))
	req.respondWith(top)
}

func statTeamActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatTeamActivity(req.login, structs.Map(stat))
	req.respondWith(activity)
}

func statUserTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatUserTop(req.login, structs.Map(stat))
	req.respondWith(top)
}

func statUserActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatUserActivity(req.login, structs.Map(stat))
	req.respondWith(activity)
}

func statRepoTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatRepoTop(req.login, structs.Map(stat))
	req.respondWith(top)
}

func statRepoActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatRepoActivity(req.login, structs.Map(stat))
	req.respondWith(activity)
}
