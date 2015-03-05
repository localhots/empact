package server

import (
	"encoding/json"
	"net/http"

	"github.com/localhots/empact/db"
)

func orgsListHandler(w http.ResponseWriter, r *http.Request) {
	jsonHandler(w, r)
	login := currentUser(r)
	orgs := db.UserOrgs(login)

	b, _ := json.Marshal(orgs)
	w.Write(b)
}
