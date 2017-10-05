package main

import (
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

import "log"

func main() {
	dsn := os.Getenv("DSN")
	workersNumber := parseIntEnv("WORKERS")
	accountsNumber := parseIntEnv("ACCOUNTS")
	transactionsNumber := parseIntEnv("TRANSACTIONS")

	config := Config{
		Dsn:                dsn,
		WorkersNumber:      workersNumber,
		AccountsNumber:     accountsNumber,
		TransactionsNumber: transactionsNumber,
	}

	dispatcher := NewDispatcher(config)
	dispatcher.Run()
}

func parseIntEnv(envName string) int {
	value, err := strconv.Atoi(os.Getenv(envName))

	if err != nil {
		log.Fatalf("Wrong value of %s variable\n", envName)
	}

	return value
}
