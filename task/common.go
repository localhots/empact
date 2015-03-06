package task

import (
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
		Quota:     int64(res.Limit),
		Remaining: int64(res.Remaining),
		ResetAt:   res.Reset.Time,
	}
	tok.Save()
}
