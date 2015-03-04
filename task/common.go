package task

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

func newGithubClient(token string) *github.Client {
	trans := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	return github.NewClient(trans.Client())
}
