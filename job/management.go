package job

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/fatih/structs"
)

type (
	Task   struct{}
	report struct {
		duration time.Duration
		err      error
	}
)

func newID() string {
	return uuid.New()
}

func (t *Task) report(rep report) {
	meta := structs.Map(t)
	meta["duration"] = rep.duration
	meta["error"] = rep.err
}
