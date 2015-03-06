package task

import (
	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func SyncRepos(token, owner string) {
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
