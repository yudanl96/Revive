package gapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/util"
	validate "github.com/yudanl96/revive/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (response *pb.CreateUserResponse, err error) {
	if violations := validateCreateUserRequest(request); violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

func validateCreateUserRequest(request *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validate.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validate.ValidateEmail(request.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}
