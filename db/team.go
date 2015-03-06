package db

type Team struct {
	ID    uint64 `json:"id"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

const orgTeamsQuery = `select * from teams where owner = ?`

func OrgTeams(login string) (teams []*Team) {
	if err := conn.Select(&teams, orgTeamsQuery, login); err != nil {
		panic(err)
	}
	return
}
