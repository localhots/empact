package db

import (
	"time"
)

type Repo struct {
	ID          uint64    `json:"id"`
	OrgID       uint64    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsFork      bool      `json:"is_fork"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Repo) Save() {
	defer measure("SaveRepo", time.Now())
	mustExecN(`
		insert into repos (owner, name, updated_at)
		values (:owner, :name, now())
		on duplicate key update
		updated_at=now()
	`, r)
}

func OrgRepos(login string) (repos []*Repo) {
	defer measure("OrgRepos", time.Now())
	mustSelect(&repos, `
		select *
		from repos r
		left join orgs o on r.org_id = o.id
		where o.login = ?
	`, login)
	return
}
