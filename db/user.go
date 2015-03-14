package db

import "time"

type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	ID        uint64 `json:"id"`
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
insert into users (login, name, id, avatar_url)
values (:login, :name, :id, :avatar_url)
on duplicate key update
login=values(login), name=values(name), avatar_url=values(avatar_url)`

func (u *User) Save() {
	defer measure("SaveUser", time.Now())
	mustExecN(saveUserQuery, u)
}

func OrgUsers(login string) (users []*User) {
	defer measure("OrgUsers", time.Now())
	mustSelect(&users, orgUsersQuery, login)
	return
}
