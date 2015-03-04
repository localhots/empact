package db

import (
	"github.com/fatih/structs"
)

type (
	Repo struct {
		Owner string
		Name  string
	}
)

const (
	saveRepoQuery = "replace into repos (owner, name, updated_at) values (?, ?, now())"
)

func (r *Repo) Save() {
	if _, err := stmt(saveRepoQuery).Exec(structs.Values(r)); err != nil {
		panic(err)
	}
}
