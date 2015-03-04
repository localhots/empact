package job

import (
	"sync"
	"time"
)

type (
	worker struct {
		id      string
		job     string
		wg      *sync.WaitGroup
		actor   func()
		reports chan<- report
		orders  <-chan order
	}
)

func (w *worker) workHard() {
	defer w.wg.Done()
	for {
		select {
		case o := <-w.orders:
			switch o {
			case stop:
				return
			default:
				panic("Confused")
			}
		default:
			action()
		}
	}
}

func (w *worker) action() {
	start := time.Now()
	defer func() {
		err := recover()
		w.reports <- report{
			duration: time.Since(start),
			success:  (err == nil),
		}
	}()

	w.actor()
}
