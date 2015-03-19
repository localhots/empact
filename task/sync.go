package task

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func SyncRepos(token, owner string) {
	defer report("SyncRepos", time.Now())
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
	defer report("SyncContrib", time.Now())
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

func SyncUserOrgs(token string) (err error) {
	defer report("SyncUserOrgs", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	for {
		var orgs []github.Organization
		var resp *github.Response
		if orgs, resp, err = client.Organizations.List("", opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, org := range orgs {
			var company, avatarURL string
			if org.Company != nil {
				company = *org.Company
			}
			if org.AvatarURL != nil {
				avatarURL = *org.AvatarURL
			}
			o := &db.Org{
				ID:        uint64(*org.ID),
				Login:     *org.Login,
				Company:   company,
				AvatarURL: avatarURL,
			}
			go SyncOrgTeams(token, o)
			go SyncOrgMembers(token, o)
			o.Save()
		}
		if opt.Page == resp.LastPage {
			break
		}
	}

	return
}
		}
	}
}
