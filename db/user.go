package db

import "time"

type User struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url" db:"avatar_url"`
}

const orgUsersQuery = `
select
    u.*
from members m
join teams t on
    m.team_id = t.id
join users u on
    m.user = u.login
where m.org = ?`

const saveUserQuery = `
insert into users (id, login, name, avatar_url, created_at, updated_at)
values (:id, :login, :name, :avatar_url, now(), now())
on duplicate key update
login = values(login),
name = values(name),
avatar_url = values(avatar_url),
updated_at = now()`

func (u *User) Save() {
	defer measure("SaveUser", time.Now())
	mustExecN(saveUserQuery, u)
}

func OrgUsers(login string) (users []*User) {
	defer measure("OrgUsers", time.Now())
	mustSelect(&users, orgUsersQuery, login)
	return
}
