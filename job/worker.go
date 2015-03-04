package job

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/localhots/steward/db"
)

type (
	worker struct {
		id       string
		job      *Job
		shutdown <-chan struct{}
	}
)

func (w *worker) workHard() {
	defer w.wg.Done()
	for {
		select {
		case <-w.shutdown:
			return
		case t := <-w.job.tasks:
			w.perform(t)
		}
	}
}

func (w *worker) perform(t *db.Task) {
	t.Worker = w.id
	t.StartedAt = time.Now()
	defer func() {
		err := recover()
		t.Duration = time.Since(t.StartedAt).Nanoseconds()
		t.Error = err.String()
		t.Save()
	}()

	w.job.actor(t)
}

func newID() string {
	return uuid.New()
}
