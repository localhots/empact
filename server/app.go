package server

import (
	"net/http"
)

func appHelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	helloTmpl.ExecuteTemplate(w, "hello", map[string]interface{}{})
}

func appAppHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	appTmpl.ExecuteTemplate(w, "app", map[string]interface{}{})
}
