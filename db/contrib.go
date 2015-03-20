package db

import (
	"time"
)

type Contrib struct {
	Week      uint64 `json:"week"`
	OrgID     uint64 `json:"org_id"`
	RepoID    uint64 `json:"repo_id"`
	UserID    uint64 `json:"user_id"`
	Commits   uint64 `json:"commits"`
	Additions uint64 `json:"additions"`
	Deletions uint64 `json:"deletions"`
}

func (c *Contrib) Save() {
	defer measure("SaveContrib", time.Now())
	mustExecN(`
		insert into contribs (week, org_id, repo_id, user_id, commits, additions, deletions)
		values (:week, :org_id, :repo_id, :user_id, :commits, :additions, :deletions)
		on duplicate key update
			commits = values(commits),
			additions = values(additions),
			deletions = values(deletions)
	`, c)
}
