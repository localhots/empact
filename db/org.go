package db

type Org struct {
	Login     string `json:"login"`
	Descr     string `json:"descr"`
	ID        uint64 `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

const userOrgsQuery = `select o.* from members m join orgs o on o.login = m.org where user = ?`

func UserOrgs(login string) (orgs []*Org) {
	mustSelect(&orgs, userOrgsQuery, login)
	return
}
