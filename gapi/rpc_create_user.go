package gapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (response *pb.CreateUserResponse, err error) {

	hashPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
	}

	arg := db.CreateUserParams{
		ID:       uuid.NewString(),
		Password: hashPassword,
		Email:    request.GetEmail(),
		Username: request.GetUsername(),
	}

	err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username or email already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to create user: %s", err)
	}

	user, err := server.store.GetUserById(ctx, arg.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to find the created user: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return res, nil
}
