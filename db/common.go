package db

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	conn  *sqlx.DB
	stmts map[string]*sqlx.Stmt
)

func Connect(params string) (err error) {
	conn, err = sqlx.Connect("mysql", params)

	conn.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	stmts = map[string]*sqlx.Stmt{}
	return
}

func stmt(query string) *sqlx.Stmt {
	if stmt, ok := stmts[query]; ok {
		return stmt
	} else {
		stmt := prepareStatement(query)
		stmts[query] = stmt
		return stmt
	}
}

func prepareStatement(query string) *sqlx.Stmt {
	if stmt, err := conn.Preparex(query); err == nil {
		return stmt
	} else {
		fmt.Println(query)
		panic(err)
	}
}
