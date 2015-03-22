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
		r     *http.Request
		w     http.ResponseWriter
		sid   string
		token string
		login string
	}
	statRequest struct {
		Org  string `structs:"org"`
		Team string `structs:"team"`
		User string `structs:"user"`
		Repo string `structs:"repo"`
		From int64  `structs:"from"`
		To   int64  `structs:"to"`
		Item string `structs:"item"`
	}
)

func parseRequest(w http.ResponseWriter, r *http.Request) (*request, *statRequest) {
	sid := sessionID(w, r)
	login, _ := redis.String(redisPool.Get().Do("HGET", "sessions", sid))
	token, _ := redis.String(redisPool.Get().Do("HGET", "tokens", sid))
	req := &request{
		r:     r,
		w:     w,
		sid:   sid,
		token: token,
		login: login,
	}
	sr := parseStatRequest(r)

	// XXX: Hack for demo account
	if req.login == "" {
		req.login = "andrewarrow"
	}

	return req, sr
}

func (r *request) authorize(token, login string) {
	redisPool.Get().Do("HSET", "tokens", r.sid, token)
	redisPool.Get().Do("HSET", "sessions", r.sid, login)
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
	} else {
		from = 0
	}
	if r.FormValue("to") != "" {
		if to, err = strconv.ParseInt(r.FormValue("to"), 10, 64); err != nil {
			panic(err)
		}
	} else {
		to = time.Now().Unix()
	}

	var item string
	switch val := r.FormValue("item"); val {
	case "author", "user":
		item = "u.login"
	case "team":
		item = "t.name"
	default:
		item = "r.name"
	}

	return &statRequest{
		Org:  r.FormValue("org"),
		Team: r.FormValue("team"),
		User: r.FormValue("user"),
		Repo: r.FormValue("repo"),
		From: from,
		To:   to,
		Item: item,
	}
}

func sessionID(w http.ResponseWriter, r *http.Request) string {
	var cook *http.Cookie
	var err error
	if cook, err = r.Cookie(cookieName); err != nil {
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
	return cook.Value
}
