package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/daniel-adam-ce/go-bank/api"
	db "github.com/daniel-adam-ce/go-bank/db/sqlc"
	"github.com/daniel-adam-ce/go-bank/gapi"
	"github.com/daniel-adam-ce/go-bank/pb"
	"github.com/daniel-adam-ce/go-bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	go runGatewayServer(config, store)
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

	log.Printf("started grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("error starting grpc server: ", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterGoBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.APIUrl)
	if err != nil {
		log.Fatal("error creating listener: ", err)
	}

	log.Printf("started gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("error starting HTTP gateway server: ", err)
	}
}
