package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SaveOrgMembers(orgID int, members []int) {
	defer measure("SaveOrgMembers", time.Now())
	tx := db.MustBegin()

	var ids = []string{}
	for _, id := range members {
		ids = append(ids, strconv.Itoa(id))
	}
	tx.MustExec(fmt.Sprintf(`
        delete from org_members
        where
            org_id = %d and
            user_id not in (%s)
    `, orgID, strings.Join(ids, ", ")))

	var values = []string{}
	for _, id := range members {
		values = append(values, fmt.Sprintf("(%d, %d)", orgID, id))
	}
	tx.MustExec(`
        insert into org_members (org_id, user_id)
        values ` + strings.Join(values, ", ") + `
        on duplicate key update org_id = org_id
    `)

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func SaveTeamMembers(orgID, teamID int, members []int) {
	defer measure("SaveTeamMembers", time.Now())
	tx := db.MustBegin()

	var ids = []string{}
	for _, id := range members {
		ids = append(ids, strconv.Itoa(id))
	}
	tx.MustExec(fmt.Sprintf(`
        delete from team_members
        where
            org_id = %d and
            team_id = %d and
            user_id not in (%s)
    `, orgID, teamID, strings.Join(ids, ", ")))

	var values = []string{}
	for _, id := range members {
		values = append(values, fmt.Sprintf("(%d, %d, %d)", orgID, teamID, id))
	}
	tx.MustExec(`
        insert into team_members (org_id, team_id, user_id)
        values ` + strings.Join(values, ", ") + `
        on duplicate key update org_id = org_id
    `)

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func SaveTeamRepos(orgID, teamID int, repos []int) {
	defer measure("SaveTeamRepos", time.Now())
	tx := db.MustBegin()

	var ids = []string{}
	for _, id := range repos {
		ids = append(ids, strconv.Itoa(id))
	}
	tx.MustExec(fmt.Sprintf(`
        delete from team_repos
        where
            org_id = %d and
            team_id = %d and
            repo_id not in (%s)
    `, orgID, teamID, strings.Join(ids, ", ")))

	var values = []string{}
	for _, id := range repos {
		values = append(values, fmt.Sprintf("(%d, %d, %d)", orgID, teamID, id))
	}
	tx.MustExec(`
        insert into team_repos (org_id, team_id, repo_id)
        values ` + strings.Join(values, ", ") + `
        on duplicate key update org_id = org_id
    `)

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
