package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
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
	var (
		path     string
		fd       *os.File
		contents []byte
		err      error
	)

	flag.StringVar(&path, "config", "config.json", "Path to configuration file")
	flag.Parse()

	if fd, err = os.Open(path); err != nil {
		panic(err)
	}
	if contents, err = ioutil.ReadAll(fd); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(contents, &conf); err != nil {
		panic(err)
	}
}
