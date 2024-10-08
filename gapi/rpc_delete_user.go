package gapi

import (
	"context"

	"github.com/yudanl96/revive/pb"
	validate "github.com/yudanl96/revive/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (response *pb.DeleteUserResponse, err error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	if violations := validateDeleteUserRequest(request); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.UserID != request.GetId() {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized for other user: %s", err)
	}

	err = server.store.DeleteUser(ctx, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to Delete user: %s", err)
	}

	res := &pb.DeleteUserResponse{}
	return res, nil
}

func validateDeleteUserRequest(request *pb.DeleteUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUUID(request.GetId()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	return violations
}
