package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	//need _ because we are not explicitly using it
)

const (
	dbDriver = "mysql"
	dbSource = "root:secret@tcp(127.0.0.1:3306)/revive?parseTime=true"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connect, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	testQueries = New(connect)

	os.Exit(m.Run())
}
