package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
	"fmt"
)

var (
	pool *sqlx.DB
	err error
)

func setupDatabase() {
	if pool, err = sqlx.Open("mysql", config.DatabaseConnectionString); err != nil {
		panic(fmt.Sprintf("Unable to open database connection. %s", err))
	}

	pool.SetMaxIdleConns(10)
	pool.SetMaxOpenConns(100)
	pool.SetConnMaxLifetime(time.Hour)

	if err = pool.Ping(); err != nil {
		panic(fmt.Sprintf("Unable to ping database. %s", err))
	}
}