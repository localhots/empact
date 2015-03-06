package server

import (
	"net/http"

	"github.com/localhots/empact/db"
)

func statOrgReposTop(w http.ResponseWriter, r *http.Request) {
	top := db.StatOrgReposTop(db.ParseRequest(r))
	respondWith(w, top)
}

func statOrgReposActivity(w http.ResponseWriter, r *http.Request) {
	activity := db.StatOrgReposActivity(db.ParseRequest(r))
	respondWith(w, activity)
}

func statOrgTeamsTop(w http.ResponseWriter, r *http.Request) {
	top := db.StatOrgTeamsTop(db.ParseRequest(r))
	respondWith(w, top)
}

func statOrgTeamsActivity(w http.ResponseWriter, r *http.Request) {
	activity := db.StatOrgTeamsActivity(db.ParseRequest(r))
	respondWith(w, activity)
}

func statOrgUsersTop(w http.ResponseWriter, r *http.Request) {
	top := db.StatOrgUsersTop(db.ParseRequest(r))
	respondWith(w, top)
}
