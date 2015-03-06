package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	db *sqlx.DB
)

func Connect(params string) (err error) {
	db, err = sqlx.Connect("mysql", params)
	db.Mapper = reflectx.NewMapper("json")
	return
}

func mustExecN(query string, arg interface{}) {
	if _, err := db.NamedExec(query, arg); err != nil {
		panic(err)
	}
}

func mustSelect(dest interface{}, query string, args ...interface{}) {
	if err := db.Select(dest, query, args...); err != nil {
		panic(err)
	}
}

func measure(op string, start time.Time) {
	duration := time.Since(start).Nanoseconds()
	outcome := "succeeded"
	if err := recover(); err != nil {
		outcome = "failed"
	}

	log.Printf("Operation %s %s; time: %d (%dms)\n", op, outcome, duration, duration/1000000)
}
