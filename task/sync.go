package task

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func SyncRepos(token, owner string) {
	report("SyncRepos", time.Now())
	client := newGithubClient(token)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{},
	}

	for {
		opt.Page++
		repos, resp, err := client.Repositories.ListByOrg(owner, opt)
		saveResponseMeta(token, resp)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			r := &db.Repo{
				Owner: owner,
				Name:  *repo.Name,
			}
			r.Save()
		}
		if len(repos) < 30 {
			break
		}
	}
}

func SyncContrib(token, owner, repo string) {
	report("SyncContrib", time.Now())
	client := newGithubClient(token)
	contribs, resp, err := client.Repositories.ListContributorsStats(owner, repo)
	saveResponseMeta(token, resp)
	if err != nil {
		if err.Error() == "EOF" {
			return // Empty repository, not an actual error
		}
		panic(err)
	}

	for _, c := range contribs {
		for _, week := range c.Weeks {
			if *week.Commits == 0 {
				continue
			}

			contrib := &db.Contrib{
				Week:      uint64(week.Week.Time.Unix()),
				Author:    *c.Author.Login,
				Owner:     owner,
				Repo:      repo,
				Commits:   uint64(*week.Commits),
				Additions: uint64(*week.Additions),
				Deletions: uint64(*week.Deletions),
			}
			contrib.Save()
		}
	}
}
