package db

import (
	"time"
)

type (
	Token struct {
		ID        int
		Owner     string
		Token     string
		Limit     int
		Remaining int
		ResetAt   time.Time
		CreatedAt time.Time
	}
)
