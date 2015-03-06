package task

import (
	"github.com/localhots/empact/db"
)

func SyncContrib(token, owner, repo string) {
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
