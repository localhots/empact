package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/GeertJohan/go.rice"
)

const (
	sessionCookie = "session_id"
)

var (
	helloTmpl = template.New("hello")
)

func init() {
	box := rice.MustFindBox("static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(box.HTTPBox())))
	tmplText, _ := box.String("hello.tmpl")
	helloTmpl, _ = helloTmpl.Parse(tmplText)

	http.HandleFunc("/", sessionHandler)
	http.HandleFunc("/auth/hello", authHelloHandler)
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
}

func Start() {
	fmt.Println("Starting server at http://localhost:8080")
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
