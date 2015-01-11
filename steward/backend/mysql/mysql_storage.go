package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/localhots/steward/steward"
)

type (
	MysqlStorage struct {
		db    *sql.DB
		State map[string]*steward.State

		importStmt    *sql.Stmt
		saveStateStmt *sql.Stmt
		loadStateStmt *sql.Stmt
	}
)

const (
	importQuery    = "replace into commits (sha1, author, repo, ts) values (?, ?, ?, ?)"
	saveStateQuery = "replace into state (repo, sha1, ts) values (?, ?, ?)"
	loadStateQuery = "select repo, sha1, ts from state"
)

func New(host, user, pass, db string) *MysqlStorage {
	var (
		s = &MysqlStorage{
			State: map[string]*steward.State{},
		}
		err         error
		databaseURI = makeDatabaseURI(host, user, pass, db)
	)

	if s.db, err = sql.Open("mysql", databaseURI); err != nil {
		panic(err)
	}
	if s.importStmt, err = s.db.Prepare(importQuery); err != nil {
		panic(err)
	}
	if s.saveStateStmt, err = s.db.Prepare(saveStateQuery); err != nil {
		panic(err)
	}

	s.loadGlobalState()

	return s
}

func (ms *MysqlStorage) Import(repo string, hist map[string]*steward.Commit) {
	var (
		lastTimestamp *time.Time
		lastSha1      string
	)

	// fmt.Println("saving", len(hist), "commits")
	// pretty.Println(hist)

	for sha1, c := range hist {
		if _, err := ms.importStmt.Exec(sha1, c.Author, repo, c.Timestamp); err != nil {
			panic(err)
		}
		if lastTimestamp == nil || lastTimestamp.After(c.Timestamp) {
			lastTimestamp = &c.Timestamp
			lastSha1 = sha1
		}
	}

	if len(hist) > 0 {
		ms.saveRepoState(repo, lastSha1, *lastTimestamp)
	}
}

func (ms *MysqlStorage) saveRepoState(repo string, sha1 string, ts time.Time) {
	ms.State[repo] = &steward.State{
		Sha1:      sha1,
		Timestamp: ts,
	}
	if _, err := ms.saveStateStmt.Exec(repo, sha1, ts); err != nil {
		panic(err)
	}
}

func (ms *MysqlStorage) loadGlobalState() {
	var (
		repo string
		sha1 string
		ts   time.Time
	)

	rows, err := ms.db.Query(loadStateQuery)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		if err := rows.Scan(&repo, &sha1, &ts); err != nil {
			panic(err)
		}
		ms.State[repo] = &steward.State{
			Sha1:      sha1,
			Timestamp: ts,
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
