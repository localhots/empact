package task

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

type (
	Tasker interface {
		Save()
		T() *db.Task
	}
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
	db.UpdateToken(&db.Token{
		Token:     token,
		Limit:     res.Limit,
		Remaining: res.Remaining,
		ResetAt:   res.Reset.Time,
	})
}
