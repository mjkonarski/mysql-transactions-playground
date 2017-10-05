package main

import "database/sql"
import "log"

type Worker struct {
	Transfers []Transfer
	Db        *sql.DB
}

func NewWorker(dsn string, transfers []Transfer) *Worker {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Unable to connect to the database with parameters: %s\n", dsn)
	}

	return &Worker{
		Transfers: transfers,
		Db:        db,
	}
}
