package server

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/localhots/empact/db"
)

func statOrgTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatOrgTop(structs.Map(stat))
	req.respondWith(top)
}

func statOrgActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatOrgActivity(structs.Map(stat))
	req.respondWith(activity)
}

func statTeamTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatTeamTop(structs.Map(stat))
	req.respondWith(top)
}

func statTeamActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatTeamActivity(structs.Map(stat))
	req.respondWith(activity)
}

func statUserTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatUserTop(structs.Map(stat))
	req.respondWith(top)
}

func statUserActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatUserActivity(structs.Map(stat))
	req.respondWith(activity)
}

func statRepoTopHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	top := db.StatRepoTop(structs.Map(stat))
	req.respondWith(top)
}

func statRepoActivityHandler(w http.ResponseWriter, r *http.Request) {
	req, stat := parseRequest(w, r)
	activity := db.StatRepoActivity(structs.Map(stat))
	req.respondWith(activity)
}
