package db

import (
	"time"
)

type Token struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Token     string    `json:"token"`
	Quota     int       `json:"quota"`
	Remaining int       `json:"remaining"`
	ResetAt   time.Time `json:"reset_at" db:"reset_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (t *Token) Save() {
	defer measure(time.Now(), "SaveToken")
	mustExecN(`
		insert into tokens (user, token, quota, remaining, reset_at, created_at)
		values (:user, :token, :quota, :remaining, :reset_at, now())
		on duplicate key update
			quota = values(quota),
			remaining = values(remaining),
			reset_at = values(reset_at)
	`, t)
}
