package gapi

import (
	"context"
	"database/sql"

	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/redisdb"
	"github.com/yudanl96/revive/util"
	validate "github.com/yudanl96/revive/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (response *pb.LoginUserResponse, err error) {
	if violations := validateLoginUserRequest(request); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := server.store.RetrieveIdByUsername(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "fail to find the user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "fail to access database: %s", err)
	}

	user, err := server.store.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "fail to find the user: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "fail to access database: %s", err)
	}

	err = util.MatchPassword(user.Password, request.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Incorrect password: %s", err)
	}

	token, payload, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create access token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create refresh token: %s", err)
	}

	mdt := server.extractMetadata(ctx)
	session := redisdb.Session{
		ID:           refreshPayload.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mdt.UserAgent, //ToDo: update this
		ClienIp:      mdt.ClientIP,
		IsBlocked:    false,
		ExpiredTime:  refreshPayload.ExpiredTime,
	}

	err = server.redisRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create session: %s", err)
	}

	res := &pb.LoginUserResponse{
		User:                    convertUser(user),
		Token:                   token,
		SessionId:               session.ID,
		RefreshToken:            refreshToken,
		TokenExpiresTime:        timestamppb.New(payload.ExpiredTime),
		RefreshTokenExpiresTime: timestamppb.New(refreshPayload.ExpiredTime),
	}

	return res, nil
}

func validateLoginUserRequest(request *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validate.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
