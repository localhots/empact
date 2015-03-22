package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/localhots/empact/config"
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
	http.HandleFunc("/api/orgs", apiUserOrgsHandler)
	http.HandleFunc("/api/teams", apiOrgTeamsHandler)
	http.HandleFunc("/api/users", apiOrgUsersHandler)
	http.HandleFunc("/api/repos", apiOrgReposHandler)
	http.HandleFunc("/api/weeks", apiOrgWeekRangeHandler)

	http.HandleFunc("/api/stat/orgs/top", statOrgTopHandler)
	http.HandleFunc("/api/stat/orgs/activity", statOrgActivityHandler)
	http.HandleFunc("/api/stat/teams/top", statTeamTopHandler)
	http.HandleFunc("/api/stat/teams/activity", statTeamActivityHandler)
	http.HandleFunc("/api/stat/users/top", statUserTopHandler)
	http.HandleFunc("/api/stat/users/activity", statUserActivityHandler)
	http.HandleFunc("/api/stat/repos/top", statRepoTopHandler)
	http.HandleFunc("/api/stat/repos/activity", statRepoActivityHandler)
}

func Start() {
	log.Printf("Starting server at http://localhost:%d\n", config.C().ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.C().ServerPort), nil)
}

func dialRedis() (redis.Conn, error) {
	return redis.Dial("tcp", fmt.Sprintf(":%d", config.C().RedisPort))
}
