package github

import (
	"fmt"
	"time"

	"code.google.com/p/goauth2/oauth"
	gh "github.com/google/go-github/github"
	"github.com/kr/pretty"
	"github.com/localhots/steward/steward"
)

const (
	DefaultPerPage = 30
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

		if len(repos) < DefaultPerPage {
			break
		}
	}
	fmt.Print("\n")

	return names
}

func (c *GithubClient) ListCommits(repo string, until *time.Time) (hist map[string]*steward.Commit, hasMore bool) {
	hist = map[string]*steward.Commit{}

	opt := &gh.CommitsListOptions{}
	if until != nil {
		opt.Until = *until
	}

	commits, _, err := c.client.Repositories.ListCommits(c.owner, repo, opt)
	if err != nil {
		fmt.Println("Error fetching commits: ", err.Error())
		return
	}

	// fmt.Println("Fetched", len(commits), "commits until", opt.Until)
	hasMore = (len(commits) == DefaultPerPage)
	for _, c := range commits {
		commit, err := makeCommit(&c)
		if err != nil {
			fmt.Println("Error:", err.Error())
			continue
		}
		hist[commit.Sha1] = commit
	}

	return
}

func makeCommit(c *gh.RepositoryCommit) (commit *steward.Commit, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print("\n\n\nTroubles with commit:")
			pretty.Println(c)
			fmt.Println("")
			panic(err)
		}
	}()

	commit = &steward.Commit{}
	if c.SHA != nil {
		commit.Sha1 = *c.SHA
	} else {
		return nil, fmt.Errorf("Missing commit SHA1 field")
	}

	if c.Author != nil {
		commit.Author = *c.Author.Login
	} else {
		return nil, fmt.Errorf("Missing author field")
	}

	if c.Commit != nil {
		if c.Commit.Author != nil {
			commit.Timestamp = *c.Commit.Author.Date
		} else {
			return nil, fmt.Errorf("Missing commit author field")
		}
	} else {
		return nil, fmt.Errorf("Missing commit field")
	}

	return
}
