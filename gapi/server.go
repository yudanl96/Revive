package gapi

import (
	"fmt"

	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/redisdb"
	"github.com/yudanl96/revive/token"
	"github.com/yudanl96/revive/util"
)

// serve grpc requests
type Server struct {
	pb.UnimplementedReviveServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
	redisRepo  *redisdb.RedisRepo
}

// creates a new grpc server
func NewServer(config util.Config, store db.Store, r *redisdb.RedisRepo) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("fail to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		redisRepo:  r,
	}

	return server, nil
}
