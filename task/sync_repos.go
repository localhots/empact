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
	repos := fetchRepos(newGithubClient(t.Token), t.Owner)
	for _, repo := range repos {
		db.ImportRepo(&db.Repo{
			Owner: t.Owner,
			Name:  repo,
		})
	}
}

func fetchRepos(client *github.Client, owner string) {
	var (
		names = []string{}
		opt   = &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{},
		}
	)

	for {
		opt.Page++
		repos, resp, err := client.Repositories.ListByOrg(owner, opt)
		// c.saveResponseMeta(resp) // Save current usage/limits
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
