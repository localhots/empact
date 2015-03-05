package db

import (
	"time"
)

type (
	Token struct {
		ID        int
		User      string
		Token     string
		Limit     int
		Remaining int
		ResetAt   time.Time
		CreatedAt time.Time
	}
)

const (
	saveTokenQuery = "" +
		"insert into tokens (`user`, token, `limit`, remaining, reset_at, created_at) " +
		"values (?, ?, ?, ?, ?, now()) " +
		"on duplicate key update " +
		"`limit` = values(`limit`), remaining = values(remaining), reset_at = values(reset_at)"
)

func (t *Token) Save() {
	if _, err := stmt(saveTokenQuery).Exec(t.User, t.Token, t.Limit, t.Remaining, t.ResetAt); err != nil {
		panic(err)
	}
}
