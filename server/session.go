package server

import (
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

const (
	sessionCookie = "session_id"
)

var (
	redisC = redis.NewPool(dialRedis, 10)
)

func dialRedis() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if cook, err := r.Cookie(sessionCookie); err != nil {
		cook = &http.Cookie{
			Name:     sessionCookie,
			Value:    uuid.New(),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cook)
	}
}

func sessionID(r *http.Request) string {
	cook, _ := r.Cookie(sessionCookie)
	return cook.Value
}

func currentUser(r *http.Request) string {
	conn := redisC.Get()
	login, _ := redis.String(conn.Do("HGET", "sessions", sessionID(r)))

	return login
}

func authorize(r *http.Request, login string) {
	redisC.Get().Do("HSET", "sessions", sessionID(r), login)
}
