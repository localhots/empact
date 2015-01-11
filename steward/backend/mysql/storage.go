package mysql

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/localhots/steward/steward"
)

type (
	Config struct {
		Hostname string
		Port     int
		Username string
		Password string
		Database string
	}
	Storage struct {
		db         *sql.DB
		importStmt *sql.Stmt
	}
)

const (
	importQuery = "" +
		"replace into contributions (author, repo, week, commits, additions, deletions) " +
		"values (?, ?, ?, ?, ?, ?)"
)

func New(c *Config) *Storage {
	var (
		s   = &Storage{}
		err error
	)

	if s.db, err = sql.Open("mysql", c.URI()); err != nil {
		panic(err)
	}
	if s.importStmt, err = s.db.Prepare(importQuery); err != nil {
		panic(err)
	}

	return s
}

func (s *Storage) ImportContributions(contrib []*steward.Contribution) {
	for _, c := range contrib {
		if _, err := s.importStmt.Exec(
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

func (c *Config) URI() (uri string) {
	if c.Username == "" && c.Password != "" {
		panic("Password is set but no username specified")
	}

	if c.Username != "" {
		uri += c.Username
	}
	if c.Password != "" {
		uri += ":" + c.Password
	}
	if c.Username != "" {
		uri += "@"
	}
	if c.Hostname != "" {
		uri += c.Hostname
	}
	if c.Port != 0 {
		uri += ":" + strconv.Itoa(c.Port)
	}
	uri += "/" + c.Database + "?parseTime=true"

	return
}
