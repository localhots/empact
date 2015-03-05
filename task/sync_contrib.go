package task

import (
	"github.com/localhots/empact/db"
)

type (
	SyncContribTask struct {
		Repo string
		*db.Task
	}
)

func SyncContrib(tk Tasker) {
	t := tk.(*SyncContribTask)
	client := newGithubClient(t.Token)
	contribs, resp, err := client.Repositories.ListContributorsStats(t.Owner, t.Repo)
	saveResponseMeta(t.Token, resp)
	if err != nil {
		if err.Error() == "EOF" {
			// Empty repository, not an actual error
			return
		}
		panic(err)
	}

	for _, c := range contribs {
		for _, week := range c.Weeks {
			if *week.Commits == 0 {
				continue
			}

			contrib := &db.Contrib{
				Week:      week.Week.Time.Unix(),
				Author:    *c.Author.Login,
				Owner:     t.Owner,
				Repo:      t.Repo,
				Commits:   *week.Commits,
				Additions: *week.Additions,
				Deletions: *week.Deletions,
			}
			contrib.Save()
		}
	}
}
