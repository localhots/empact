package db

import (
	"time"
)

type Org struct {
	GithubID  uint64 `json:"github_id"`
	Login     string `json:"login"`
	Company   string `json:"company"`
	AvatarURL string `json:"avatar_url"`
}

const userOrgsQuery = `select o.* from members m join orgs o on o.login = m.org where user = ?`

func UserOrgs(login string) (orgs []*Org) {
	defer measure("UserOrgs", time.Now())
	mustSelect(&orgs, userOrgsQuery, login)
	return
}
