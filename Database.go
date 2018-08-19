package main

import (
	db "github.com/gocraft/dbr"
	_ "github.com/go-sql-driver/mysql"
)

var (
	pool *db.Connection
	err error
)

func setupDatabase() {
	pool, err = db.Open("mysql", config.DatabaseConnectionString, nil)
	if err != nil {
		panic(err)
	}

	pool.SetMaxOpenConns(100)
}