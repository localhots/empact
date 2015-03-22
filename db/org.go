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

func OrgWeekRange(login string) (min int, max int) {
	row := db.QueryRow(`
		select
			min(c.week) as min,
			max(c.week) as max
		from contribs c
		join orgs o on o.id = c.org_id
		where o.login = ?
	`, login)
	if err := row.Scan(&min, &max); err != nil {
		panic(err)
	}
	return
}
