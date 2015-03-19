package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SaveOrgMembers(org uint64, members []uint64) {
	defer measure("SaveOrgMembers", time.Now())
	tx := db.MustBegin()

	var ids = []string{}
	for _, id := range members {
		ids = append(ids, strconv.FormatUint(id, 10))
	}
	tx.MustExec(fmt.Sprintf(`
        delete from org_members
        where
            org_id = %d and
            user_id not in (%s)
    `, org, strings.Join(ids, ", ")))

	var values = []string{}
	for _, id := range members {
		values = append(values, fmt.Sprintf("(%d, %d)", org, id))
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
