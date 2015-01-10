package steward

import (
	"time"
)

type (
	State struct {
		Sha1      string
		Timestamp time.Time
	}
)
