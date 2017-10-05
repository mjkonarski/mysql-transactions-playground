package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "os"
import "log"
import "fmt"

func main() {
  dsn := os.Getenv("DSN")
  db, err := sql.Open("mysql", dsn)

  if err != nil {
    log.Fatal("Unable to connect to the database with parameters: %s\n", dsn)
  }

  rows, err := db.Query("SELECT id, balance FROM accounts")
  if err != nil {
	   log.Fatal(err)
  }

  for rows.Next() {
  	var id int32
    var balance int32

  	if err := rows.Scan(&id, &balance); err != nil {
  		log.Fatal(err)
  	}

  	fmt.Printf("%d is %d\n", id, balance)
  }

  if err := rows.Err(); err != nil {
  	log.Fatal(err)
  }

  err = db.Close()
  if err != nil {
    log.Fatal(err)
  }
}
