package task

import (
	"fmt"
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
	t := &db.Token{
		Token:     token,
		Quota:     res.Limit,
		Remaining: res.Remaining,
		ResetAt:   res.Reset.Time,
	}
	db.Later(func() { t.Save() })
}

func report(start time.Time, format string, args ...interface{}) {
	duration := time.Since(start).Nanoseconds()
	err := recover()
	outcome := "succeeded"
	if err != nil {
		outcome = "failed"
		defer panic(err)
	}

	task := fmt.Sprintf(format, args...)
	log.Printf("Task %s %s; time: %d (%dms)\n", task, outcome, duration, duration/1000000)
}
