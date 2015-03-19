package db

import (
	"time"
)

type Org struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	Company   string `json:"company"`
	AvatarURL string `json:"avatar_url" db:"avatar_url"`
}

const userOrgsQuery = `select o.* from members m join orgs o on o.login = m.org where user = ?`
const saveOrgQuery = `
insert into orgs (id, login, company, avatar_url, created_at, updated_at)
values (:id, :login, :company, :avatar_url, now(), now())
on duplicate key update
login = values(login),
company = values(company),
avatar_url = values(avatar_url),
updated_at=now()`

func (o *Org) Save() {
	defer measure("SaveOrg", time.Now())
	mustExecN(saveOrgQuery, o)
}

func UserOrgs(login string) (orgs []*Org) {
	defer measure("UserOrgs", time.Now())
	mustSelect(&orgs, userOrgsQuery, login)
	return
}
