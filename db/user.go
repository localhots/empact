package db

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Save() {
	defer measure("SaveUser", time.Now())
	mustExecN(`
		insert into users (id, login, name, avatar_url, created_at, updated_at)
		values (:id, :login, :name, :avatar_url, now(), now())
		on duplicate key update
			login = values(login),
			name = values(name),
			avatar_url = values(avatar_url),
			updated_at = now()
	`, u)
}

func OrgUsers(login string) (users []*User) {
	defer measure("OrgUsers", time.Now())
	mustSelect(&users, `
		select u.*
		from org_members om
		join users u on u.id = om.user_id
		join orgs o on o.id = om.org_id
		where o.login = ?
	`, login)
	return
}
