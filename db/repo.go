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
	repoImportQuery = "replace into repos (owner, name, updated_at) values (?, ?, now())"
)

func ImportRepo(r *Repo) {
	if _, err := stmt(repoImportQuery).Exec(structs.Values(r)); err != nil {
		panic(err)
	}
}
