package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/GeertJohan/go.rice"
)

var (
	helloTmpl = template.New("hello")
	appTmpl   = template.New("app")
	box       = rice.MustFindBox("app")
)

func init() {
	parseTemplate("hello.tmpl", helloTmpl)
	parseTemplate("app.tmpl", appTmpl)

	// Serving static files
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(box.HTTPBox())))

	http.HandleFunc("/", sessionHandler)
	http.HandleFunc("/hello", appHelloHandler)
	http.HandleFunc("/app", appAppHandler)
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
	http.HandleFunc("/api/", authHandler)
	http.HandleFunc("/api/orgs", apiOrgsHandler)
	http.HandleFunc("/api/teams", apiTeamsHandler)
	http.HandleFunc("/api/repos", apiReposHandler)
}

func Start() {
	fmt.Println("Starting server at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func parseTemplate(file string, tmpl *template.Template) {
	if tmplText, err := box.String(file); err == nil {
		tmpl, _ = tmpl.Parse(tmplText)
	} else {
		panic(err)
	}
}

func respondWith(w http.ResponseWriter, resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Write(b)
}
