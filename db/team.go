package db

import (
	"time"
)

type Team struct {
	ID    uint64 `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

const orgTeamsQuery = `select * from teams where owner = ?`

func OrgTeams(login string) (teams []*Team) {
	defer measure("OrgTeams", time.Now())
	mustSelect(&teams, orgTeamsQuery, login)
	return
}
