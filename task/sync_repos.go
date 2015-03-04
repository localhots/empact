package task

import (
	"github.com/google/go-github/github"
	"github.com/localhots/steward/db"
	"github.com/localhots/steward/job"
)

type (
	SyncReposTask struct {
		Owner string
		Token string
		job.Task
	}
)

func SyncRepos(t SyncReposTask) {
	repos := fetchRepos(t.Token, t.Owner)
	for _, repo := range repos {
		db.ImportRepo(&db.Repo{
			Owner: t.Owner,
			Name:  repo,
		})
	}
}

func fetchRepos(token, owner string) {
	var (
		client = newGithubClient(token)
		names  = []string{}
		opt    = &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{},
		}
	)

	for {
		opt.Page++
		repos, resp, err := client.Repositories.ListByOrg(owner, opt)
		saveResponseMeta(token, resp)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			names = append(names, *repo.Name)
		}
		if len(repos) < 30 {
			break
		}
	}

	return names
}
