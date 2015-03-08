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
