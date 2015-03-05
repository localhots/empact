package server

import (
	"net/http"
)

func Start() {
	http.HandleFunc("/auth/signin", authSigninHandler)
	http.HandleFunc("/auth/callback", authCallbackHandler)
	http.ListenAndServe(":8080", nil)
}
