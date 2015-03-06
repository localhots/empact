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
	conn.Mappper = reflectx.NewMapperFunc("json", strings.ToLower)
	return
}
