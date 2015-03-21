package db

import (
	"time"
)

type Org struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Company   string    `json:"company"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (o *Org) Save() {
	defer measure(time.Now(), "SaveOrg")
	mustExecN(`
		insert into orgs (id, login, company, avatar_url, updated_at)
		values (:id, :login, :company, :avatar_url, now())
		on duplicate key update
			login = values(login),
			company = values(company),
			avatar_url = values(avatar_url),
			updated_at = now()
	`, o)
}

func UserOrgs(login string) (orgs []*Org) {
	defer measure(time.Now(), "UserOrgs")
	mustSelect(&orgs, `
		select o.*
		from org_members m
		join users u on u.id = m.user_id
		join orgs o on o.id = m.org_id
		where u.login = ?
	`, login)
	return
}
