package job

import (
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/fatih/structs"
	"github.com/localhots/empact/task"
)

type (
	Job struct {
		Name    string
		actor   func(task.Tasker)
		workers map[string]*worker
		orders  chan struct{} // Currently shutdown only
		tasks   chan task.Tasker
		wg      sync.WaitGroup
	}
)

var (
	jobs = map[string]*Job{}
)

func Enqueue(t task.Tasker) {
	dt := t.T()
	dt.Job = structs.Name(t)
	dt.CreatedAt = time.Now()

	j, ok := jobs[dt.Job]
	if !ok {
		switch dt.Job {
		case "FetchAccessTokenTask":
			j = New(dt.Job, task.FetchAccessToken)
		case "SyncContribTask":
			j = New(dt.Job, task.SyncContrib)
		case "SyncReposTask":
			j = New(dt.Job, task.SyncRepos)
		default:
			panic("Unknown task: " + dt.Job)
		}
		jobs[dt.Job] = j
		j.Resize(1)
	}
	j.tasks <- t
}

func New(name string, actor func(task.Tasker)) *Job {
	return &Job{
		Name:    name,
		actor:   actor,
		workers: make(map[string]*worker),
		orders:  make(chan struct{}),
		tasks:   make(chan task.Tasker, 1000),
	}
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
				id:  uuid.New(),
				job: j,
			}
			go w.workHard()
			j.workers[w.id] = w
		}
		j.wg.Add(del)
	} else {
		for i := 0; i > del; i-- {
			j.orders <- struct{}{}
		}
	}
}
