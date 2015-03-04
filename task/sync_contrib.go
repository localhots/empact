package task

import (
	"github.com/localhots/steward/db"
	"github.com/localhots/steward/job"
)

type (
	SyncContribTask struct {
		Owner string
		Repo  string
		Token string
		job.Task
	}
)

func SyncContrib(t SyncContribTask) {
	contribs := fetchContrib(t.Token, t.Owner, t.Repo)
	for _, c := range contribs {
		db.ImportRepo(c)
	}
}

func fetchContrib(token, owner, repo string) (res []*db.Contrib) {
	client := newGithubClient(token)
	contribs, resp, err := client.Repositories.ListContributorsStats(owner, repo)
	saveResponseMeta(token, resp)
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

			res = append(res, &db.Contrib{
				Week:      week.Week.Time.Unix(),
				Author:    *c.Author.Login,
				Owner:     owner,
				Repo:      repo,
				Commits:   *week.Commits,
				Additions: *week.Additions,
				Deletions: *week.Deletions,
			})
		}
	}

	return
}
