package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"github.com/yudanl96/revive/chat"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/gapi"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/redisdb"
	"github.com/yudanl96/revive/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Unable to load configuration: ", err)
	}

	var wg sync.WaitGroup
	client, err := startRedisServer(config, &wg)
	if err != nil {
		log.Fatal("Failed to create Redis server: ", err)
	}
	//make sure redis is initialized before continuing
	wg.Wait()

	r := redisdb.RedisRepo{
		Client: client,
	}

	connect, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := db.NewStore(connect)

	mux := http.NewServeMux()

	newroom := chat.NewRoom()

	go newroom.Run()
	mux.Handle("/room", newroom)

	go startGrpcServer(config, store, &r)
	startGWServer(mux, config, store, &r)
}

func startRedisServer(config util.Config, wg *sync.WaitGroup) (*redis.Client, error) {
	wg.Add(1)
	defer wg.Done() // Notify that the Redis server setup is complete

	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisServerAddress,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("Failed to connect to redis: ", err)
		return nil, err
	}

	log.Printf("start redis server at %s", rdb.Options().Addr)
	return rdb, nil
}

func startGrpcServer(config util.Config, store db.Store, r *redisdb.RedisRepo) {
	server, err := gapi.NewServer(config, store, r)
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

func startGWServer(mux *http.ServeMux, config util.Config, store db.Store, r *redisdb.RedisRepo) {
	server, err := gapi.NewServer(config, store, r)
	if err != nil {
		log.Fatal("Failed to create server: ", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //will be executed before exiting function
	//prevent system doing unnesseary work

	err = pb.RegisterReviveHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("Failed to register handler server: ", err)
	}

	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Failed to create lisenter: ", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Failed to start HTTP gateway server: ", err)
	}
}

// func startGinServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal("Failed to create server: ", err)
// 	}

// 	err = server.Start(config.HTTPServerAddress)
// 	if err != nil {
// 		log.Fatal("Failed to connect to server: ", err)
// 	}
// }
