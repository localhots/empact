package db

import (
	"github.com/fatih/structs"
)

type (
	Contrib struct {
		Week      int64
		Author    string
		Owner     string
		Repo      string
		Commits   int
		Additions int
		Deletions int
	}
)

const (
	importContribQuery = "" +
		"replace into contributions (week, author, owner, repo, commits, additions, deletions) " +
		"values (?, ?, ?, ?, ?, ?, ?)"
)

func ImportContrib(c *Contrib) {
	if _, err := stmt(importContribQuery).Exec(structs.Values(c)); err != nil {
		panic(err)
	}
}
