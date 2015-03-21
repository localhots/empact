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
		ID:        *user.ID,
		AvatarURL: avatarURL,
	}

	return
}

func SyncUserInfo(token, login string) (err error) {
	defer report("SyncUserInfo", time.Now())
	var u *db.User
	if u, err = FetchUserInfo(token, login); err == nil {
		db.Queue(func() { u.Save() })
	}
	return
}
