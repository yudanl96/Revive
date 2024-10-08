package gapi

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yudanl96/revive/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewToken(ctx context.Context, request *pb.RenewTokenRequest) (response *pb.RenewTokenResponse, err error) {
	//if violations := validateRenewTokenRequest(request); violations != nil {
	//	return nil, invalidArgumentError(violations)
	//}

	refreshPayload, err := server.tokenMaker.VerifyToken(request.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token: %s", err)
	}

	session, err := server.redisRepo.RetrieveSession(ctx, refreshPayload.ID.String())
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.NotFound, "fail to find session: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "fail to access database: %s", err)
	}

	if session.IsBlocked {
		return nil, status.Errorf(codes.Unauthenticated, "blocked session: %s", err)
	}

	if session.Username != refreshPayload.Username {
		return nil, status.Errorf(codes.Unauthenticated, "wrong session user: %s", err)
	}

	if session.RefreshToken != request.RefreshToken {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token mismatch: %s", err)
	}

	token, payload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create access token: %s", err)
	}

	res := &pb.RenewTokenResponse{
		Token: &pb.Token{
			Token:            token,
			TokenExpiresTime: timestamppb.New(payload.ExpiredTime)},
	}

	return res, nil
}

// func validateRenewTokenRequest(request *pb.RenewTokenRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if err := validate.ValidateUsername(request.GetUsername()); err != nil {
// 		violations = append(violations, fieldViolation("username", err))
// 	}
// 	if err := validate.ValidatePassword(request.GetPassword()); err != nil {
// 		violations = append(violations, fieldViolation("password", err))
// 	}
// 	return violations
// }
