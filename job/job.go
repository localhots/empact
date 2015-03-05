package job

import (
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/localhots/empact/db"
)

type (
	Job struct {
		Name    string
		actor   func(*db.Task)
		workers map[string]*worker
		tasks   chan *db.Task
		wg      sync.WaitGroup
	}
)

func New(name string, actor func()) *Job {
	&Job{
		Name:    name,
		actor:   actor,
		workers: make(map[string]*worker),
		tasks:   make(chan *db.Task, 1000),
	}
}

func (j *Job) Perform(t *db.Task) {
	t.Job = structs.Name(t)
	t.CreatedAt = time.Now()
	j.tasks <- t
}

func (j *Job) Size() int {
	return len(j.workers)
}

func (j *Job) Stop() {
	j.Resize(0)
	j.wg.Wait()
}

func (j *Job) Resize(n int) {
	if n < 0 {
		n = 0
	}
	if del := n - len(j.workers); del > 0 {
		for i := 0; i < del; i++ {
			w := &worker{
				id:       newID(),
				job:      j,
				shutdown: make(<-chan struct{}, 1),
			}
			go w.workHard()
			j.workers[w.id] = w
		}
		j.wg.Add(del)
	} else {
		for i := 0; i > del; i-- {
			j.shutdown <- struct{}{}
		}
	}
}
