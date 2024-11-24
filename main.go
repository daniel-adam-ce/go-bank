package main

import (
	"database/sql"
	"log"

	"github.com/daniel-adam-ce/go-bank/api"
	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/daniel-adam-ce/go-bank/util"
	_ "github.com/lib/pq"
)

func main() {

	var err error

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// conn, err := pgxpool.New(context.Background(), dbSource)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.APIUrl)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
