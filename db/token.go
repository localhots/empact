package db

import (
	"time"
)

type Token struct {
	ID        uint64    `json:"id"`
	User      string    `json:"user"`
	Token     string    `json:"token"`
	Quota     uint64    `json:"quota"`
	Remaining uint64    `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
	CreatedAt time.Time `json:"created_at"`
}

const saveTokenQuery = `
insert into tokens (user, token, quota, remaining, reset_at, created_at)
values (:user, :token, :quota, :remaining, :reset_at, now())
on duplicate key update
quota = values(quota), remaining = values(remaining), reset_at = values(reset_at)`

func (t *Token) Save() {
	mustExecN(saveTokenQuery, t)
}
