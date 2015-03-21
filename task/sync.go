package task

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func SyncOrgRepos(token string, org *db.Org) {
	defer report(time.Now(), "SyncOrgRepos (%s)", org.Login)
	client := newGithubClient(token)
	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {
		opt.Page++
		var repos []github.Repository
		var resp *github.Response
		var err error
		if repos, resp, err = client.Repositories.ListByOrg(org.Login, opt); err != nil {
			panic(err)
		}
		saveResponseMeta(token, resp)

		for _, repo := range repos {
			var descr string
			if repo.Description != nil {
				descr = *repo.Description
			}
			r := &db.Repo{
				ID:          *repo.ID,
				OrgID:       *repo.Owner.ID,
				Name:        *repo.Name,
				Description: descr,
				IsPrivate:   *repo.Private,
				IsFork:      *repo.Fork,
			}
			go SyncContrib(token, org.Login, r)
			db.Queue(func() { r.Save() })
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
}

func SyncContrib(token, owner string, repo *db.Repo) {
	defer report(time.Now(), "SyncContrib (%s/%s)", owner, repo.Name)
	client := newGithubClient(token)

	var contribs []github.ContributorStats
	var resp *github.Response
	var err error
	if contribs, resp, err = client.Repositories.ListContributorsStats(owner, repo.Name); err != nil {
		if err.Error() == "EOF" {
			return // Empty repository, not an actual error
		}
		panic(err)
	}
	saveResponseMeta(token, resp)

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

func SyncUserOrgs(token string) {
	defer report(time.Now(), "SyncUserOrgs")
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	for {
		opt.Page++
		var orgs []github.Organization
		var resp *github.Response
		var err error
		if orgs, resp, err = client.Organizations.List("", opt); err != nil {
			panic(err)
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
			go SyncOrgRepos(token, o)
			db.Queue(func() { o.Save() })
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}

	return
}

func SyncOrgTeams(token string, org *db.Org) {
	defer report(time.Now(), "SyncOrgTeams (%s)", org.Login)
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	for {
		opt.Page++
		var teams []github.Team
		var resp *github.Response
		var err error
		if teams, resp, err = client.Organizations.ListTeams(org.Login, opt); err != nil {
			panic(err)
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

func SyncOrgMembers(token string, org *db.Org) {
	defer report(time.Now(), "SyncOrgMembers (%s)", org.Login)
	client := newGithubClient(token)
	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	var ids = []int{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		var err error
		if users, resp, err = client.Organizations.ListMembers(org.Login, opt); err != nil {
			panic(err)
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

func SyncTeamMembers(token string, team *db.Team) {
	defer report(time.Now(), "SyncTeamMembers (%d/%s)", team.OrgID, team.Name)
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	var ids = []int{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		var err error
		if users, resp, err = client.Organizations.ListTeamMembers(int(team.ID), opt); err != nil {
			panic(err)
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

func SyncTeamRepos(token string, team *db.Team) {
	defer report(time.Now(), "SyncTeamRepos (%d/%s)", team.OrgID, team.Name)
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	var ids = []int{}
	for {
		opt.Page++
		var repos []github.Repository
		var resp *github.Response
		var err error
		if repos, resp, err = client.Organizations.ListTeamRepos(int(team.ID), opt); err != nil {
			panic(err)
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
