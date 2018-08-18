package main

import (
	db "github.com/gocraft/dbr"
	_ "github.com/go-sql-driver/mysql"
)

const CONNECTION_STRING = "root:DEluxe8892@/runemanager"
var (
	pool *db.Connection
	err error
)

func setupDatabase() {
	pool, err = db.Open("mysql", CONNECTION_STRING, nil)
	if err != nil {
		panic(err)
	}

	pool.SetMaxOpenConns(100)
}