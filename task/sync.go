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
			r := &db.Repo{
				Owner: owner,
				Name:  *repo.Name,
			}
			r.Save()
		}
		if opt.Page >= resp.LastPage {
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

	for _, contrib := range contribs {
		for _, week := range contrib.Weeks {
			if *week.Commits == 0 {
				continue
			}

			c := &db.Contrib{
				Week:      uint64(week.Week.Time.Unix()),
				Author:    *contrib.Author.Login,
				Owner:     owner,
				Repo:      repo,
				Commits:   uint64(*week.Commits),
				Additions: uint64(*week.Additions),
				Deletions: uint64(*week.Deletions),
			}
			c.Save()
		}
	}
}

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
				ID:         uint64(*team.ID),
				Name:       *team.Name,
				Slug:       *team.Slug,
				Permission: *team.Permission,
				OrgID:      org.ID,
			}
			go SyncTeamMembers(token, t)
			t.Save()
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

	var ids = []uint64{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		if users, resp, err = client.Organizations.ListMembers(org.Login, opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, user := range users {
			var name, avatarURL string
			if user.Name != nil {
				name = *user.Name
			}
			if user.AvatarURL != nil {
				avatarURL = *user.AvatarURL
			}
			u := &db.User{
				ID:        uint64(*user.ID),
				Login:     *user.Login,
				Name:      name,
				AvatarURL: avatarURL,
			}
			u.Save()
			ids = append(ids, u.ID)
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
	db.SaveOrgMembers(org.ID, ids)

	return
}

func SyncTeamMembers(token string, team *db.Team) (err error) {
	defer report("SyncTeamMembers", time.Now())
	client := newGithubClient(token)
	opt := &github.ListOptions{PerPage: 100}

	var ids = []uint64{}
	for {
		opt.Page++
		var users []github.User
		var resp *github.Response
		if users, resp, err = client.Organizations.ListTeamMembers(int(team.ID), opt); err != nil {
			return
		}
		saveResponseMeta(token, resp)

		for _, user := range users {
			ids = append(ids, uint64(*user.ID))
		}
		if opt.Page >= resp.LastPage {
			break
		}
	}
	db.SaveTeamMembers(team.OrgID, team.ID, ids)

	return
}
