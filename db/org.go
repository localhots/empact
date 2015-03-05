package db

type (
	Org struct {
		Login     string
		Descr     string
		ID        int64
		AvatarURL string
	}
)

const (
	userOrgsQuery = "select org from members where user = ?"
)

func UserOrgs(login string) (orgs []string) {
	if res, err := stmt(userOrgsQuery).Query(login); err == nil {
		defer res.Close()
		for res.Next() {
			var org string
			if err := res.Scan(&org); err != nil {
				panic(err)
			}
			orgs = append(orgs, org)
		}
	} else {
		panic(err)
	}
	return
}
