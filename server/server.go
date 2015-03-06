package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
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
	http.HandleFunc("/api/", authHandler)
	http.HandleFunc("/api/orgs", apiOrgsHandler)
	http.HandleFunc("/api/teams", apiTeamsHandler)
	http.HandleFunc("/api/repos", apiReposHandler)
	http.HandleFunc("/auth/hello", authHelloHandler)
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
}

func Start() {
	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
