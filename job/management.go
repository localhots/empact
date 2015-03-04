package job

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

type (
	order  byte
	report struct {
		duration time.Duration
		success  bool
	}
)

const (
	stop order = iota
)

func newID() string {
	return uuid.New()
}
