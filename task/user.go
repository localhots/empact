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

	u = &db.User{
		Login:     *user.Login,
		Name:      *user.Name,
		ID:        uint64(*user.ID),
		AvatarURL: *user.AvatarURL,
	}

	return
}
