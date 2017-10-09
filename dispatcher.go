package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const MAX_Balance = 1000

type Config struct {
	Dsn             string
	WorkersNumber   int
	AccountsNumber  int
	TransfersNumber int
}

type Dispatcher struct {
	Config  Config
	Db      *sql.DB
	Workers []*Worker
}

func NewDispatcher(config Config) *Dispatcher {
	return &Dispatcher{
		Config: config,
	}
}

func (dispatcher *Dispatcher) Run() {
	dispatcher.connectToDb()
	dispatcher.prepareAccounts()

	transfers := dispatcher.generateTransfers()
	// generatedSum := dispatcher.sumTransfers(transfers)
	checkingWorker := NewCheckingWorker(dispatcher.Config.Dsn)

	checkingChannel := make(chan bool)
	go checkingWorker.Run(checkingChannel)
	dispatcher.runWorkers(transfers)
	close(checkingChannel)
	fetchedSum := dispatcher.fetchSumFromDB()

	if fetchedSum == 0 {
		fmt.Println("Sum in the DB is correct!")
	} else {
		fmt.Println("Sum in the DB is incorrect!")
	}
}

func (dispatcher *Dispatcher) connectToDb() {
	db, err := sql.Open("mysql", dispatcher.Config.Dsn)

	if err != nil {
		log.Fatalf("Unable to connect to the DB: %s", dispatcher.Config.Dsn)
	} else {
		log.Print("Successfully connected to the db")
	}

	dispatcher.Db = db
}

func (dispatcher *Dispatcher) prepareAccounts() {
	_, err := dispatcher.Db.Exec("TRUNCATE accounts")

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := dispatcher.Db.
		Prepare("INSERT INTO accounts(id, balance) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for accountID := 1; accountID <= dispatcher.Config.AccountsNumber; accountID++ {
		_, err := stmt.Exec(accountID, 0)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Added %d accounts", dispatcher.Config.AccountsNumber)
}

func (dispatcher *Dispatcher) generateTransfers() []Transfer {
	transfers := make([]Transfer, dispatcher.Config.TransfersNumber)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < dispatcher.Config.TransfersNumber; i++ {
		from := r.Intn(dispatcher.Config.AccountsNumber) + 1
		to := r.Intn(dispatcher.Config.AccountsNumber) + 1

		for to == from {
			to = r.Intn(dispatcher.Config.AccountsNumber) + 1
		}

		Balance := r.Intn(MAX_Balance) + 1
		transfer := Transfer{
			From:    from,
			To:      to,
			Balance: Balance,
		}

		transfers[i] = transfer
	}

	log.Printf("Created %d trasfers", len(transfers))

	return transfers
}

func (dispatcher *Dispatcher) sumTransfers(transfers []Transfer) map[int]int {
	sums := make(map[int]int)

	for accountID := 1; accountID <= dispatcher.Config.AccountsNumber; accountID++ {
		sums[accountID] = 0
	}

	for _, transfer := range transfers {
		sums[transfer.From] -= transfer.Balance
		sums[transfer.To] += transfer.Balance
	}

	return sums
}

func (dispatcher *Dispatcher) runWorkers(transfers []Transfer) {
	dispatcher.Workers = make([]*Worker, dispatcher.Config.WorkersNumber)
	workerTransfers := make(map[int][]Transfer)
	for workerID := 0; workerID <= dispatcher.Config.WorkersNumber; workerID++ {
		workerTransfers[workerID] = make([]Transfer, 0)
	}

	for transferID, transfer := range transfers {
		workerID := transferID % dispatcher.Config.WorkersNumber
		workerTransfers[workerID] = append(workerTransfers[workerID], transfer)
	}

	workerFinishChannel := make(chan bool, dispatcher.Config.WorkersNumber)

	for workerID := 0; workerID < dispatcher.Config.WorkersNumber; workerID++ {
		worker := NewWorker(dispatcher.Config.Dsn)
		go worker.Run(workerTransfers[workerID], workerFinishChannel)
	}

	log.Printf("Ran %d workers", dispatcher.Config.WorkersNumber)

	for workerID := 0; workerID < dispatcher.Config.WorkersNumber; workerID++ {
		<-workerFinishChannel
	}
}

func (dispatcher *Dispatcher) fetchSumFromDB() int {
	var sum int
	err := dispatcher.Db.QueryRow("SELECT SUM(balance) FROM accounts").
		Scan(&sum)
	if err != nil {
		log.Fatal(err)
	}

	return sum
}
