package task

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func SyncRepos(token, owner string) {
	defer report("SyncRepos", time.Now())
	client := newGithubClient(token)
	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {
		opt.Page++
		repos, resp, err := client.Repositories.ListByOrg(owner, opt)
		saveResponseMeta(token, resp)
		if err != nil {
			panic(err)
		}
		for _, repo := range repos {
			var descr string
			if repo.Description != nil {
				descr = *repo.Description
			}
			r := &db.Repo{
				OrgID:       *repo.Organization.ID,
				Name:        *repo.Name,
				Description: descr,
				IsPrivate:   *repo.Private,
				IsFork:      *repo.Fork,
			}
			db.Queue(func() { r.Save() })
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
}

func SyncContrib(token, owner string, repo *db.Repo) {
	defer report("SyncContrib", time.Now())
	client := newGithubClient(token)

	contribs, resp, err := client.Repositories.ListContributorsStats(owner, repo.Name)
	saveResponseMeta(token, resp)
	if err != nil {
		if err.Error() == "EOF" {
			return // Empty repository, not an actual error
		}
		panic(err)
	}

	for _, contrib := range contribs {
		for _, week := range contrib.Weeks {
			if *week.Commits == 0 {
				continue
			}

			c := &db.Contrib{
				Week:      int(week.Week.Time.Unix()),
				OrgID:     repo.OrgID,
				RepoID:    repo.ID,
				UserID:    *contrib.Author.ID,
				Commits:   *week.Commits,
				Additions: *week.Additions,
				Deletions: *week.Deletions,
			}
			db.Queue(func() { c.Save() })
		}
	}
}

func SyncUserOrgs(token string) (err error) {
	defer report("SyncUserOrgs", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	for {
		opt.Page++
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
				ID:        *org.ID,
				Login:     *org.Login,
				Company:   company,
				AvatarURL: avatarURL,
			}
			go SyncOrgTeams(token, o)
			go SyncOrgMembers(token, o)
			db.Queue(func() { o.Save() })
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}

	return
}

func SyncOrgTeams(token string, org *db.Org) (err error) {
	defer report("SyncOrgTeams", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	for {
		opt.Page++
		var teams []github.Team
		var resp *github.Response
		if teams, resp, err = client.Organizations.ListTeams(org.Login, opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, team := range teams {
			t := &db.Team{
				Name:       *team.Name,
				ID:         *team.ID,
				Slug:       *team.Slug,
				Permission: *team.Permission,
				OrgID:      org.ID,
			}
			go SyncTeamMembers(token, t)
			go SyncTeamRepos(token, t)
			db.Queue(func() { t.Save() })
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}

	return
}

func SyncOrgMembers(token string, org *db.Org) (err error) {
	defer report("SyncOrgMembers", time.Now())
	client := newGithubClient(token)
	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	var ids = []int{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		if users, resp, err = client.Organizations.ListMembers(org.Login, opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, user := range users {
			ids = append(ids, *user.ID)
			go SyncUserInfo(token, *user.Login)
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
	db.Queue(func() { db.SaveOrgMembers(org.ID, ids) })

	return
}

func SyncTeamMembers(token string, team *db.Team) (err error) {
	defer report("SyncTeamMembers", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	var ids = []int{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		if users, resp, err = client.Organizations.ListTeamMembers(int(team.ID), opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, user := range users {
			ids = append(ids, *user.ID)
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
	db.Queue(func() { db.SaveTeamMembers(team.OrgID, team.ID, ids) })

	return
}

func SyncTeamRepos(token string, team *db.Team) (err error) {
	defer report("SyncTeamRepos", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	var ids = []int{}
	for {
		opt.Page++
		var repos []github.Repository
		var resp *github.Response
		if repos, resp, err = client.Organizations.ListTeamRepos(int(team.ID), opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, repo := range repos {
			ids = append(ids, *repo.ID)
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
	db.Queue(func() { db.SaveTeamRepos(team.OrgID, team.ID, ids) })

	return
}
