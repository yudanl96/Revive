package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yudanl96/revive/api"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load configuration: ", err)
	}
	fmt.Println(config.DBDriver, config.DBSourse, config.ServerAddress)
	connect, err := sql.Open(config.DBDriver, config.DBSourse)

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(connect)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to connect to server: ", err)
	}
}
