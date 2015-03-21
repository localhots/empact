package db

import (
	"time"
)

type Contrib struct {
	Week      int `json:"week"`
	OrgID     int `json:"org_id" db:"org_id"`
	RepoID    int `json:"repo_id" db:"repo_id"`
	UserID    int `json:"user_id" db:"user_id"`
	Commits   int `json:"commits"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

func (c *Contrib) Save() {
	defer measure(time.Now(), "SaveContrib")
	mustExecN(`
		insert into contribs (week, org_id, repo_id, user_id, commits, additions, deletions)
		values (:week, :org_id, :repo_id, :user_id, :commits, :additions, :deletions)
		on duplicate key update
			commits = values(commits),
			additions = values(additions),
			deletions = values(deletions)
	`, c)
}
