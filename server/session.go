package server

import (
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

const (
	cookieName = "session_id"
)

var (
	redisPool = redis.NewPool(dialRedis, 10)
)

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if cook, err := r.Cookie(cookieName); err != nil {
		cook = &http.Cookie{
			Name:     cookieName,
			Value:    uuid.New(),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cook)
	}
}

func createSession(r *http.Request, login string) {
	redisPool.Get().Do("HSET", "sessions", sessionID(r), login)
}

func sessionID(r *http.Request) string {
	cook, _ := r.Cookie(cookieName)
	return cook.Value
}

func sessionUser(r *http.Request) (login string) {
	login, _ = redis.String(redisPool.Get().Do("HGET", "sessions", sessionID(r)))
	return
}

func dialRedis() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}
