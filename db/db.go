package db

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

var (
	db         *sqlx.DB
	queryQueue = make(chan func())
)

func Connect(params string) (err error) {
	db, err = sqlx.Connect("mysql", params)
	db.Mapper = reflectx.NewMapper("json")
	go processQueue()
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

func mustSelectN(dest interface{}, query string, params interface{}) {
	var stmt *sqlx.NamedStmt
	var err error
	if stmt, err = db.PrepareNamed(query); err != nil {
		panic(err)
	}
	if err = stmt.Select(dest, params); err != nil {
		panic(err)
	}
}

func Queue(fun func()) {
	queryQueue <- fun
}

func processQueue() {
	for {
		(<-queryQueue)()
	}
}

func measure(start time.Time, op string) {
	duration := time.Since(start).Nanoseconds()
	outcome := "succeeded"
	err := recover()
	if err != nil {
		outcome = "failed"
		defer panic(err)
	}

	log.Printf("Operation %s %s; time: %d (%dms)\n", op, outcome, duration, duration/1000000)
}
