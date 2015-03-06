package db

import (
	"time"
)

type Contrib struct {
	Week      uint64 `json:"week"`
	Author    string `json:"author"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Commits   uint64 `json:"commits"`
	Additions uint64 `json:"additions"`
	Deletions uint64 `json:"deletions"`
}

const saveContribQuery = `
insert into contributions (week, author, owner, repo, commits, additions, deletions)
values (:week, :author, :owner, :repo, :commits, :additions, :deletions)
on duplicate key update
commits=values(commits), additions=values(additions), deletions=values(deletions)`

func (c *Contrib) Save() {
	defer measure("SaveContrib", time.Now())
	mustExecN(saveContribQuery, c)
}
