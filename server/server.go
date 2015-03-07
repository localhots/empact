package server

import (
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
)

const (
	cookieName = "session_id"
)

var (
	redisPool = redis.NewPool(dialRedis, 10)
)

func init() {
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

func dialRedis() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}
