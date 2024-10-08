package gapi

import (
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func convertPost(post db.Post) *pb.Post {
	return &pb.Post{
		Id:          post.ID,
		UserId:      post.UserID,
		Description: post.Description,
		Price:       post.Price,
		Sold:        post.Sold,
		CreatedAt:   timestamppb.New(post.CreatedAt),
		UpdatedAt:   timestamppb.New(post.UpdatedAt),
	}
}
