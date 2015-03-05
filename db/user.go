package db

type (
	User struct {
		Login     string `json:"login"`
		Name      string `json:"name"`
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
	}
)

const (
	saveUserQuery = "" +
		"insert into users (login, name, id, avatar_url) " +
		"values (?, ?, ?, ?) " +
		"on duplicate key update " +
		"login=values(login), name=values(name), avatar_url=values(avatar_url)"
)

func (u *User) Save() {
	if _, err := stmt(saveUserQuery).Exec(u.Login, u.Name, u.ID, u.AvatarURL); err != nil {
		panic(err)
	}
}
