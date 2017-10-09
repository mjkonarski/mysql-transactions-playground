package main

import "database/sql"
import "log"

type Worker struct {
	Transfers []Transfer
	Db        *sql.DB
}

func NewWorker(dsn string) *Worker {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Unable to connect to the database with parameters: %s\n", dsn)
	}

	return &Worker{
		Db: db,
	}
}

func (worker *Worker) Run(transfers []Transfer, finishChannel chan<- bool) {
	updateStmt, err := worker.Db.
		Prepare("UPDATE accounts SET balance = balance + ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	for _, transfer := range transfers {
		_, err := updateStmt.Exec(-transfer.Balance, transfer.From)
		if err != nil {
			log.Fatal(err)
		}

		_, err = updateStmt.Exec(transfer.Balance, transfer.To)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Worker finished running %d transfers", len(transfers))

	finishChannel <- true
}
