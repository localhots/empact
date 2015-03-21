package db

import (
	"time"
)

type Team struct {
	ID         int       `json:"id"`
	OrgID      int       `json:"org_id" db:"org_id"`
	Slug       string    `json:"slug"`
	Name       string    `json:"name"`
	Permission string    `json:"permission"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (t *Team) Save() {
	defer measure(time.Now(), "SaveTeam")
	mustExecN(`
		insert into teams (id, org_id, slug, name, permission, updated_at)
		values (:id, :org_id, :slug, :name, :permission, now())
		on duplicate key update
			slug = values(slug),
			name = values(name),
			permission = values(permission),
			updated_at = now()
	`, t)
}

func OrgTeams(login string) (teams []*Team) {
	defer measure(time.Now(), "OrgTeams")
	mustSelect(&teams, `
		select t.*
		from teams t
		join orgs o on o.id = t.org_id
		where o.login = ?
	`, login)
	return
}
