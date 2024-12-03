package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/daniel-adam-ce/go-bank/api"
	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/daniel-adam-ce/go-bank/gapi"
	"github.com/daniel-adam-ce/go-bank/pb"
	"github.com/daniel-adam-ce/go-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	var err error

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	// fmt.Printf("conn: %s", config.DBSource)
	// conn, err := pgxpool.New(context.Background(), dbSource)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	// runGinServer(config, store)

	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	err = server.Start(config.APIUrl)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GAPIUrl)
	if err != nil {
		log.Fatal("error starting grpc server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("error starting grpc server: ", err)
	}
}
