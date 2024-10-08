package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"github.com/yudanl96/revive/util"
	validate "github.com/yudanl96/revive/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (response *pb.CreatePostResponse, err error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	if violations := validateCreatePostRequest(request); violations != nil {
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

	var description string
	if request.Genai {
		description, err = util.GenerateText(request.GetDescription(), 50)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate desciption text: %s", err)
		}
		if len(description) > 250 {
			description = description[:250]
		}
	} else {
		description = request.GetDescription()
	}

	arg := db.CreatePostParams{
		ID:          uuid.NewString(),
		UserID:      id,
		Description: description,
		Price:       request.GetPrice(),
	}

	err = server.store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to Update user: %s", err)
	}

	post, err := server.store.GetPostById(ctx, arg.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to find the Updated user: %s", err)
	}

	res := &pb.CreatePostResponse{
		Post: convertPost(post),
	}
	return res, nil
}

func validateCreatePostRequest(request *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validate.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validate.ValidateDescription(request.GetDescription()); err != nil {
		violations = append(violations, fieldViolation("description", err))
	}

	if err := validate.ValidatePrice(request.GetPrice()); err != nil {
		violations = append(violations, fieldViolation("price", err))
	}

	return violations
}
