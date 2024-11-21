package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const dbDriver = "postgres"

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbSource := os.Getenv("DB_SOURCE")
	log.Print("test", dbSource)
	// conn, err := pgxpool.New(context.Background(), dbSource)

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
