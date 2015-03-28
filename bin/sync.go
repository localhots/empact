package main

import (
	"flag"
	"fmt"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
	"github.com/localhots/empact/task"
)

func main() {
	var token string
	flag.StringVar(&token, "token", "", "GitHub access token")
	flag.Parse()
	config.Load()
	if token == "" {
		fmt.Println("Access token is required")
		return
	}

	if err := db.Connect(config.C().DatabaseURI); err != nil {
		panic(err)
	}

	task.SyncUserOrgs(token)
	select {}
}
