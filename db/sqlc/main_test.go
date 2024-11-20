package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbSource := os.Getenv("DB_SOURCE")
	log.Print("test", dbSource)
	conn, err := pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
