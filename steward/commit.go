package steward

import (
	"time"
)

type (
	Commit struct {
		Repo      string
		Sha1      string
		Author    string
		Timestamp time.Time
	}
)
