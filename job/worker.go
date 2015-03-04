package job

import (
	"time"
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

func (w *worker) perform(t Task) {
	start := time.Now()
	defer func() {
		err := recover()
		t.report(report{
			duration: time.Since(start),
			success:  err,
		})
	}()

	w.job.actor(t)
}
