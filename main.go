package main

import (
	"bank/api"
	db "bank/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/bankdb?sslmode=disable"
	address  = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Can't connect to db", err.Error())
	}

	store := db.NewSQLStore(conn)
	server := api.NewServer(&store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Can't start server", err.Error())
	}
}
