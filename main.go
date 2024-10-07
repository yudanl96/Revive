package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yudanl96/revive/api"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/gapi"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load configuration: ", err)
	}

	connect, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(connect)
	startGrpcServer(config, store)
}

func startGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Failed to connect to server: ", err)
	}
}

func startGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterReviveServer(grpcServer, server)
	reflection.Register(grpcServer) //for client to explore what is in the server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Failed to create lisenter: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Failed to start grpc server: ", err)
	}
}
