package db

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	conn *sqlx.DB
)

func Connect(params string) (err error) {
	conn, err = sqlx.Connect("mysql", params)
	conn.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return
}

func mustExecN(query string, arg interface{}) {
	if _, err := conn.NamedExec(query, arg); err != nil {
		panic(err)
	}
}

func mustSelect(dest interface{}, query string, args ...interface{}) {
	if err := conn.Select(dest, query, args...); err != nil {
		panic(err)
	}
}
