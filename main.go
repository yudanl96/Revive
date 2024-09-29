package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yudanl96/revive/api"
	db "github.com/yudanl96/revive/db/sqlc"
)

const (
	dbDriver      = "mysql"
	dbSource      = "root:secret@tcp(127.0.0.1:3306)/revive?parseTime=true"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connect, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(connect)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Failed to connect to server: ", err)
	}
}
