package task

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/localhots/empact/db"
)

func FetchUserInfo(token, login string) (u *db.User, err error) {
	defer report("FetchUserInfo", time.Now())

	client := newGithubClient(token)
	var user *github.User
	var resp *github.Response
	if user, resp, err = client.Users.Get(login); err != nil {
		return
	}
	saveResponseMeta(token, resp)

	name := ""
	if n := user.Name; n != nil {
		name = *user.Name
	}

	avatarURL := ""
	if url := user.AvatarURL; url != nil {
		avatarURL = *user.AvatarURL
	}

	u = &db.User{
		Login:     *user.Login,
		Name:      name,
		ID:        uint64(*user.ID),
		AvatarURL: avatarURL,
	}

	return
}

func FetchUserOrgs(token string) (orgs []*db.Org, err error) {
	client := newGithubClient(token)
	var ghorgs []github.Organization
	var resp *github.Response
	if ghorgs, resp, err = client.Organizations.List("", nil); err != nil {
		return
	}
	saveResponseMeta(token, resp)

	for _, ghorg := range ghorgs {
		org := &db.Org{
			GithubID:  uint64(*ghorg.ID),
			Login:     *ghorg.Login,
			Company:   *ghorg.Company,
			AvatarURL: *ghorg.AvatarURL,
		}
		orgs = append(orgs, org)
	}

	return
}
