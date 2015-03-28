package main

import (
	"flag"

	"github.com/localhots/empact/config"
	"github.com/localhots/empact/db"
	"github.com/localhots/empact/server"
)

func main() {
	flag.Parse()
	config.Load()
	if err := db.Connect(config.C().DatabaseURI); err != nil {
		panic(err)
	}

	server.Start()
}
