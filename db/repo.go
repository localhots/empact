package db

import (
	"time"
)

type Repo struct {
	ID        uint64    `json:"id"`
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPrivate bool      `json:"is_private"`
	IsForm    bool      `json:"is_fork"`
}

const orgReposQuery = `select * from repos where owner = ?`
const saveRepoQuery = `
insert into repos (owner, name, updated_at)
values (:owner, :name, now())
on duplicate key update
updated_at=now()`

func (r *Repo) Save() {
	defer measure("SaveRepo", time.Now())
	mustExecN(saveRepoQuery, r)
}

func OrgRepos(login string) (repos []*Repo) {
	defer measure("OrgRepos", time.Now())
	mustSelect(&repos, orgReposQuery, login)
	return
}
