package task

import (
	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

type (
	SyncReposTask struct {
		*db.Task
	}
)

func SyncRepos(tk Tasker) {
	t := tk.(*SyncReposTask)
	client := newGithubClient(t.Token)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{},
	}

	for {
		opt.Page++
		repos, resp, err := client.Repositories.ListByOrg(t.Owner, opt)
		saveResponseMeta(t.Token, resp)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			r := &db.Repo{
				Owner: t.Owner,
				Name:  *repo.Name,
			}
			r.Save()
		}
		if len(repos) < 30 {
			break
		}
	}
}
