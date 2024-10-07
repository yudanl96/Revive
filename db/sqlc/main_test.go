package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yudanl96/revive/util"
	//need _ because we are not explicitly using it
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Unable to load configuration: ", err)
	}

	connect, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(connect)

	os.Exit(m.Run())
}
