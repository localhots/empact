package github

import (
	"fmt"

	"code.google.com/p/goauth2/oauth"
	gh "github.com/google/go-github/github"
	"github.com/localhots/steward/steward"
)

const (
	DEFAULT_PER_PAGE = 30
)

type (
	GithubClient struct {
		owner  string
		client *gh.Client
	}
)

func New(clientID, clientSecret, owner string) *GithubClient {
	// Auth here
	return nil
}

// Temp method
func NewClient(token, owner string) *GithubClient {
	trans := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	return &GithubClient{
		owner:  owner,
		client: gh.NewClient(trans.Client()),
	}
}

func (c *GithubClient) ListRepos() []string {
	var (
		names = []string{}
		opt   = &gh.RepositoryListByOrgOptions{
			ListOptions: gh.ListOptions{},
		}
	)

	fmt.Print("Loading repositories ")
	for {
		fmt.Print(".")
		opt.Page++
		repos, _, err := c.client.Repositories.ListByOrg(c.owner, opt)
		if err != nil {
			panic(err)
		}

		for _, repo := range repos {
			names = append(names, *repo.Name)
		}

		if len(repos) < DEFAULT_PER_PAGE {
			break
		}
	}
	fmt.Print("\n")

	return names
}

func (c *GithubClient) ListCommits(repo string) map[string]*steward.Commit {
	var (
		history = map[string]*steward.Commit{}
		opt     = &gh.CommitsListOptions{}
	)

	fmt.Print(repo, " ")
	for {
		fmt.Print(".")
		commits, _, err := c.client.Repositories.ListCommits(c.owner, repo, opt)
		if err != nil {
			panic(err)
		}

		for _, c := range commits {
			history[*c.SHA] = &steward.Commit{
				Repo:      repo,
				Author:    *c.Author.Login,
				Timestamp: *c.Commit.Author.Date,
			}
			opt.Until = *c.Commit.Author.Date
		}

		if len(commits) < DEFAULT_PER_PAGE {
			break
		}
	}
	fmt.Print("\n")

	return history
}
