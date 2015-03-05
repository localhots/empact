package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	conn  *sql.DB
	stmts map[string]*sql.Stmt
)

func Connect(uri string) (err error) {
	conn, err = sql.Open("mysql", uri)
	stmts = map[string]*sql.Stmt{}
	return
}

func stmt(query string) *sql.Stmt {
	if stmt, ok := stmts[query]; ok {
		return stmt
	} else {
		stmt := prepareStatement(query)
		stmts[query] = stmt
		return stmt
	}
}

func prepareStatement(query string) *sql.Stmt {
	if stmt, err := conn.Prepare(query); err == nil {
		return stmt
	} else {
		fmt.Println(query)
		panic(err)
	}
}
