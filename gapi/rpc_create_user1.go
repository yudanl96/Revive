package gapi

import (
	"context"

	"github.com/yudanl96/revive/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (response *pb.CreateUserResponse, err error) {
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}
