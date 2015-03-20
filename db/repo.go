package db

import (
	"time"
)

type Repo struct {
	ID          int       `json:"id"`
	OrgID       int       `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
	IsFork      bool      `json:"is_fork"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Repo) Save() {
	defer measure("SaveRepo", time.Now())
	mustExecN(`
		insert into repos (org_id, name, description, is_private, is_fork, updated_at)
		values (:org_id, :name, :description, :is_private, :is_fork, now())
		on duplicate key update
			org_id = values(org_id),
			name = values(name),
			description = values(description),
			is_private = values(is_private),
			is_fork = values(is_fork),
			updated_at = now()
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
