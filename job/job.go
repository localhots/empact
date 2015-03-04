package job

import (
	"sync"
)

type (
	Job struct {
		Name    string
		actor   func()
		workers map[string]*worker
		orders  chan order
		wg      sync.WaitGroup
	}
)

func New(name string, actor func()) *Job {
	&Job{
		Name:    name,
		actor:   actor,
		workers: make(map[string]*worker),
		orders:  make(chan order),
	}
}

func (j *Job) Workers(n int) {
	if n < 0 {
		n = 0
	}
	if del := n - len(j.workers); del > 0 {
		for i := 0; i < del; i++ {
			w := &worker{
				id:    newID(),
				job:   j.Name,
				wg:    j.wg,
				actor: j.actor,
			}
			go w.workHard()
		}
		j.wg.Add(del)
	} else {
		for i := 0; i > del; i-- {
			j.orders <- stop
		}
	}
}
