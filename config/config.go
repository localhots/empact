package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type (
	Config struct {
		Domain         string `json:"app_domain"`
		DatabaseURI    string `json:"database_uri"`
		AuthURL        string `json:"github_auth_url"`
		AccessTokenURL string `json:"github_access_token_url"`
		ClientID       string `json:"github_client_id"`
		ClientSecret   string `json:"github_client_secret"`
		RedirectURI    string `json:"github_redirect_uri"`
	}
)

var (
	conf Config
)

// Config is immutable and is always returned by value
func C() Config {
	return conf
}

func init() {
	var err error
	var path string
	flag.StringVar(&path, "config", "config.json", "Path to configuration file")
	flag.Parse()

	var fd *os.File
	if fd, err = os.Open(path); err != nil {
		panic(err)
	}
	var contents []byte
	if contents, err = ioutil.ReadAll(fd); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(contents, &conf); err != nil {
		panic(err)
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime | log.Lshortfile)
}
