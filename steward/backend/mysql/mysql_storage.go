package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/localhots/steward/steward"
)

type (
	MysqlStorage struct {
		db         *sql.DB
		importStmt *sql.Stmt
	}
)

const (
	importQuery = "" +
		"replace into contributions (author, repo, week, commits, additions, deletions) " +
		"values (?, ?, ?, ?, ?, ?)"
)

func New(host, user, pass, db string) *MysqlStorage {
	var (
		s           = &MysqlStorage{}
		err         error
		databaseURI = makeDatabaseURI(host, user, pass, db)
	)

	if s.db, err = sql.Open("mysql", databaseURI); err != nil {
		panic(err)
	}
	if s.importStmt, err = s.db.Prepare(importQuery); err != nil {
		panic(err)
	}

	return s
}

func (ms *MysqlStorage) ImportContributions(contrib []*steward.Contribution) {
	for _, c := range contrib {
		if _, err := ms.importStmt.Exec(
			c.Author,
			c.Repo,
			c.Week,
			c.Commits,
			c.Additions,
			c.Deletions,
		); err != nil {
			panic(err)
		}
	}
}

func makeDatabaseURI(host, user, pass, db string) string {
	var (
		databaseURI string
	)

	if user == "" && pass != "" {
		panic("Password is set but no is user specified")
	}

	if user != "" {
		databaseURI += user
	}
	if pass != "" {
		databaseURI += ":" + pass
	}
	if user != "" {
		databaseURI += "@"
	}
	if host != "" {
		databaseURI += host
	}
	databaseURI += "/" + db + "?parseTime=true"

	return databaseURI
}
