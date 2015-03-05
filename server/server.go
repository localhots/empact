package server

import (
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

const (
	sessionCookie = "session_id"
)

func Start() {
	fmt.Println("Starting server at http://localhost:8080")
	http.HandleFunc("/", sessionHandler)
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
	http.ListenAndServe(":8080", nil)
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
