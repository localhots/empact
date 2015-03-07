package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", sessionHandler)
	http.HandleFunc("/hello", appHelloHandler)
	http.HandleFunc("/app", appAppHandler)
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
	http.HandleFunc("/api/", authHandler)
	http.HandleFunc("/api/orgs", apiOrgsHandler)
	http.HandleFunc("/api/teams", apiTeamsHandler)
	http.HandleFunc("/api/repos", apiReposHandler)
	http.HandleFunc("/api/stat/repos/top", statOrgReposTop)
	http.HandleFunc("/api/stat/repos/activity", statOrgReposActivity)
	http.HandleFunc("/api/stat/teams/top", statOrgTeamsTop)
	http.HandleFunc("/api/stat/teams/activity", statOrgTeamsActivity)
	http.HandleFunc("/api/stat/users/top", statOrgUsersTop)
}

func Start() {
	log.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func respondWith(w http.ResponseWriter, resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(b)
}
