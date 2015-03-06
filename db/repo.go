package db

type Repo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

const saveRepoQuery = `replace into repos (owner, name, updated_at) values (?, ?, now())`

func (r *Repo) Save() {
	conn.MustExec(saveRepoQuery, r.Owner, r.Name)
}
