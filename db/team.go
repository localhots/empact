package db

import (
	"time"
)

type Team struct {
	ID         uint64 `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Permission string `json:"permission"`
	OrgID      uint64 `json:"org_id" db:"org_id"`
}

const orgTeamsQuery = `select * from teams where owner = ?`
const saveTeamQuery = `
insert into teams (id, slug, name, permission, org_id, created_at, updated_at)
values (:id, :slug, :name, :permission, :org_id, now(), now())
on duplicate key update
slug = values(slug),
permission = values(permission),
updated_at=now()`

func (t *Team) Save() {
	defer measure("SaveTeam", time.Now())
	mustExecN(saveTeamQuery, t)
}

func OrgTeams(login string) (teams []*Team) {
	defer measure("OrgTeams", time.Now())
	mustSelect(&teams, orgTeamsQuery, login)
	return
}
