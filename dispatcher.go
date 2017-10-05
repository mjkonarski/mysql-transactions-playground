package main

import "math/rand"
import "time"
import "database/sql"
import "log"

const MAX_AMOUNT = 1000

type Config struct {
	Dsn                string
	WorkersNumber      int
	AccountsNumber     int
	TransactionsNumber int
}

type Dispatcher struct {
	Config Config
	Db     *sql.DB
}

func NewDispatcher(config Config) *Dispatcher {
	return &Dispatcher{
		Config: config,
	}
}

func (dispatcher *Dispatcher) Run() {
	dispatcher.connectToDb()
	dispatcher.PrepareAccounts()
}

func (dispatcher *Dispatcher) connectToDb() {
	db, err := sql.Open("mysql", dispatcher.Config.Dsn)

	if err != nil {
		log.Fatal("Unable to connect to the DB: %s", dispatcher.Config.Dsn)
	} else {
		log.Print("Successfully connected to the db")
	}

	dispatcher.Db = db
}

func (dispatcher *Dispatcher) PrepareAccounts() {
	_, err := dispatcher.Db.Exec("TRUNCATE accounts")

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := dispatcher.Db.
		Prepare("INSERT INTO accounts(id, balance) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for accountId := 1; accountId <= dispatcher.Config.AccountsNumber; accountId++ {
		_, err := stmt.Exec(accountId, 0)

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Added %d accounts", dispatcher.Config.AccountsNumber)
}

func (dispatcher *Dispatcher) generateTransfers() []Transfer {
	transfers := make([]Transfer, dispatcher.Config.TransactionsNumber)

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := 0; i < dispatcher.Config.TransactionsNumber; i++ {
		from := r.Intn(dispatcher.Config.AccountsNumber) + 1
		to := r.Intn(dispatcher.Config.AccountsNumber) + 1

		for to == from {
			to = r.Intn(dispatcher.Config.AccountsNumber) + 1
		}

		amount := r.Intn(MAX_AMOUNT) + 1
		transfer := Transfer{
			From:   from,
			To:     to,
			Amount: amount,
		}

		transfers = append(transfers, transfer)
	}

	return transfers
}

// db, err := sql.Open("mysql", dsn)
//
//
//
// rows, err := db.Query("SELECT id, balance FROM accounts")
// if err != nil {
//    log.Fatal(err)
// }
//
// for rows.Next() {
//   var id int32
//   var balance int32
//
//   if err := rows.Scan(&id, &balance); err != nil {
//     log.Fatal(err)
//   }
//
//   fmt.Printf("%d is %d\n", id, balance)
// }
//
// if err := rows.Err(); err != nil {
//   log.Fatal(err)
// }
//
// err = db.Close()
// if err != nil {
//   log.Fatal(err)
// }
