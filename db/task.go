package db

import (
	"time"
)

type (
	Task struct {
		Token     string
		Owner     string
		Job       string
		Worker    string
		Duration  int64
		Error     string
		CreatedAt time.Time
		StartedAt time.Time
	}
)

const (
	saveTaskQuery = "" +
		"insert into tasks (token, owner, job, worker, duration, error, created_at, started_at) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?)"
)

func (t *Task) Save() {
	if _, err := stmt(saveTaskQuery).Exec(
		t.Token, t.Owner,
		t.Job, t.Worker,
		t.Duration, t.Error,
		t.CreatedAt, t.StartedAt,
	); err != nil {
		panic(err)
	}
}

func (t *Task) T() *Task {
	return t
}
