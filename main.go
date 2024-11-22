package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/daniel-adam-ce/go-bank/api"
	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

func main() {

	var err error
	err = godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbSource := os.Getenv("DB_SOURCE")
	serverAddress := os.Getenv("API_URL")
	log.Print("test", dbSource)
	// conn, err := pgxpool.New(context.Background(), dbSource)

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
