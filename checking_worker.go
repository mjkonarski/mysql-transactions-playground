package main

import (
	"database/sql"
	"fmt"
	"log"
)

type CheckingWorker struct {
	Db *sql.DB
}

func NewCheckingWorker(dsn string) *CheckingWorker {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Unable to connect to the database with parameters: %s\n", dsn)
	}

	return &CheckingWorker{
		Db: db,
	}
}

func (worker *CheckingWorker) Run(finishChannel <-chan bool) {
	finish := false

	for !finish {
		var sum int
		err := worker.Db.QueryRow("SELECT SUM(balance) FROM accounts").
			Scan(&sum)

		if sum != 0 {
			fmt.Println("Sum not correct!")
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		select {
		case _, ok := <-finishChannel:
			if !ok {
				return
			}
		default:
		}
	}
}
