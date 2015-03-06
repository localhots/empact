package task

import (
	"log"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func newGithubClient(token string) *github.Client {
	trans := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	return github.NewClient(trans.Client())
}

func saveResponseMeta(token string, res *github.Response) {
	if res == nil {
		return
	}
	tok := &db.Token{
		Token:     token,
		Quota:     uint64(res.Limit),
		Remaining: uint64(res.Remaining),
		ResetAt:   res.Reset.Time,
	}
	tok.Save()
}

func report(task string, start time.Time) {
	duration := time.Since(start).Nanoseconds()
	outcome := "done"
	if err := recover(); err != nil {
		outcome = "failed"
	}

	log.Printf("Task %s %s; time: %d (%dms)\n", task, outcome, duration, duration/1000000)
}
