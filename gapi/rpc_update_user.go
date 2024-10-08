package gapi

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/util"
	validate "github.com/yudanl96/revive/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (response *pb.UpdateUserResponse, err error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	if violations := validateUpdateUserRequest(request); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != request.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized for other user: %s", err)
	}

	id, err := server.store.RetrieveIdByUsername(ctx, request.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find the user to update: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot find the user id: %s", err)
	}

	arg := db.UpdateUserParams{
		ID: id,
		Email: sql.NullString{
			String: request.GetEmail(),
			Valid:  request.Email != nil,
		},
		Username: sql.NullString{
			String: request.GetNewUsername(),
			Valid:  request.NewUsername != nil,
		},
	}

	if request.Password != nil {
		hashPassword, err := util.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
		}
		arg.Password = sql.NullString{
			String: hashPassword,
			Valid:  true,
		}
	}

	err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username or email already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to Update user: %s", err)
	}

	user, err := server.store.GetUserById(ctx, arg.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to find the Updated user: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return res, nil
}

func validateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if request.NewUsername != nil {
		if err := validate.ValidateUsername(request.GetNewUsername()); err != nil {
			violations = append(violations, fieldViolation("new username", err))
		}
	}

	if request.Password != nil {
		if err := validate.ValidatePassword(request.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if request.Email != nil {
		if err := validate.ValidateEmail(request.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	return violations
}
