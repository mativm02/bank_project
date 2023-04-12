package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mativm02/bank_system/util"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

// TestMain is the entry point for all tests.
// It sets up the database connection and runs the tests.
// It also tears down the database connection after all tests are done.
// This function is called by the testing package automatically.
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
