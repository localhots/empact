package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/garyburd/redigo/redis"
)

type (
	request struct {
		r         *http.Request
		w         http.ResponseWriter
		sessionID string
		login     string
	}
	statRequest struct {
		org  string
		team string
		user string
		from int64
		to   int64
	}
)

func parseRequest(w http.ResponseWriter, r *http.Request) (*request, *statRequest) {
	setCookie(w, r)
	cook, _ := r.Cookie(cookieName)
	login, _ := redis.String(redisPool.Get().Do("HGET", "sessions", cook.Value))
	req := &request{
		r:         r,
		w:         w,
		sessionID: cook.Value,
		login:     login,
	}
	return req, parseStatRequest(r)
}

func (r *request) authorize(login string) {
	redisPool.Get().Do("HSET", "sessions", r.sessionID, login)
}

func (r *request) respondWith(resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	r.w.Header().Set("Access-Control-Allow-Origin", "*")
	r.w.Header().Set("Content-Type", "application/json; charset=utf8")
	r.w.Write(b)
}

func parseStatRequest(r *http.Request) *statRequest {
	var err error
	var from, to int64
	if r.FormValue("from") != "" {
		if from, err = strconv.ParseInt(r.FormValue("from"), 10, 64); err != nil {
			panic(err)
		}
	}
	if r.FormValue("to") != "" {
		if to, err = strconv.ParseInt(r.FormValue("to"), 10, 64); err != nil {
			panic(err)
		}
	}
	return &statRequest{
		org:  r.FormValue("org"),
		team: r.FormValue("team"),
		user: r.FormValue("user"),
		from: from,
		to:   to,
	}
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	if cook, err := r.Cookie(cookieName); err != nil {
		cook = &http.Cookie{
			Name:     cookieName,
			Value:    uuid.New(),
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cook)
		r.AddCookie(cook)
	}
}
