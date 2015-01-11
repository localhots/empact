package github

import (
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	gh "github.com/google/go-github/github"
	"github.com/localhots/steward/steward"
)

const (
	DefaultPerPage = 30
)

type (
	GithubClient struct {
		owner     string
		client    *gh.Client
		limit     int
		remaining int
		limitEnds time.Time
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

	for {
		opt.Page++
		repos, _, err := c.client.Repositories.ListByOrg(c.owner, opt)
		if err != nil {
			panic(err)
		}

		for _, repo := range repos {
			names = append(names, *repo.Name)
		}

		if len(repos) < DefaultPerPage {
			break
		}
	}

	return names
}

func (c *GithubClient) ListContributors(repo string) []*steward.Contribution {
	var (
		contrib = []*steward.Contribution{}
	)

	cslist, resp, err := c.client.Repositories.ListContributorsStats(c.owner, repo)
	c.saveResponseMeta(resp)
	if err != nil {
		if err.Error() == "EOF" {
			// Empty repository, not an actual error
			return contrib
		}

		fmt.Println("Error loading contributors stats for repo", repo)
		fmt.Println(err.Error())
		return contrib
	}

	for _, cs := range cslist {
		for _, week := range cs.Weeks {
			if *week.Commits == 0 {
				continue
			}

			contrib = append(contrib, &steward.Contribution{
				Author:    *cs.Author.Login,
				Repo:      repo,
				Week:      week.Week.Time.Unix(),
				Commits:   *week.Commits,
				Additions: *week.Additions,
				Deletions: *week.Deletions,
			})
		}
	}

	return contrib
}

func (c *GithubClient) saveResponseMeta(res *gh.Response) {
	c.limit = res.Limit
	c.remaining = res.Remaining
	c.limitEnds = res.Reset.Time
}
