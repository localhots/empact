package main

import (
	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
	"github.com/localhots/empact/server"
)

func main() {
	if err := db.Connect(config.C().DatabaseURI); err != nil {
		panic(err)
	}

	server.Start()
}
