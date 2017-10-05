# mysql-transactions-playground
Go programs to test different transactions and locking concepts in MySQL

## Running

```
mysql < schema.sql
go build && go install
DSN="root@/transactions_playground" WORKERS=10 ACCOUNTS=2 TRANSACTIONS=1000 mysql-transactions-playground
```
