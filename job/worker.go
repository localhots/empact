package job

import (
	"sync"
	"time"

	"github.com/localhots/empact/task"
)

type (
	worker struct {
		id  string
		job *Job
		wg  *sync.WaitGroup
	}
)

func (w *worker) workHard() {
	defer w.wg.Done()
	for {
		select {
		case <-w.job.orders:
			return
		case t := <-w.job.tasks:
			w.perform(t)
		}
	}
}

func (w *worker) perform(t task.Tasker) {
	dt := t.T()
	dt.Worker = w.id
	dt.StartedAt = time.Now()
	defer func() {
		if err := recover(); err != nil {
			// dt.Error = err.(string)
		}
		dt.Duration = time.Since(dt.StartedAt).Nanoseconds()
		dt.Save()
	}()

	w.job.actor(t)
}
