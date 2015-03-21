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
		Quota:     res.Limit,
		Remaining: res.Remaining,
		ResetAt:   res.Reset.Time,
	}
	db.Queue(func() { tok.Save() })
}

func report(task string, start time.Time) {
	duration := time.Since(start).Nanoseconds()
	outcome := "succeeded"
	if err := recover(); err != nil {
		outcome = "failed"
	}

	log.Printf("Task %s %s; time: %d (%dms)\n", task, outcome, duration, duration/1000000)
}
