package db

type Contrib struct {
	Week      int64  `json:"week"`
	Author    string `json:"author"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Commits   int64  `json:"commits"`
	Additions int64  `json:"additions"`
	Deletions int64  `json:"deletions"`
}

const saveContribQuery = `
insert into contributions (week, author, owner, repo, commits, additions, deletions)
values (?, ?, ?, ?, ?, ?, ?)
on duplicate key update
commits=values(commits), additions=values(additions), deletions=values(deletions)`

func (c *Contrib) Save() {
	conn.MustExec(saveContribQuery, c.Week, c.Author, c.Owner, c.Repo, c.Commits, c.Additions, c.Deletions)
}
