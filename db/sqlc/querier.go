// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CountSavedPostsByPost(ctx context.Context, postID string) (int64, error)
	CreatePost(ctx context.Context, arg CreatePostParams) error
	CreateSavedPost(ctx context.Context, arg CreateSavedPostParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeletePost(ctx context.Context, id string) error
	DeleteSavedPost(ctx context.Context, arg DeleteSavedPostParams) error
	DeleteUser(ctx context.Context, id string) error
	GetPostById(ctx context.Context, id string) (Post, error)
	GetSavedPostByIds(ctx context.Context, arg GetSavedPostByIdsParams) (SavedPost, error)
	GetUserById(ctx context.Context, id string) (User, error)
	ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error)
	ListSavedPostsByUser(ctx context.Context, userID string) ([]SavedPost, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	RetrieveIdByEmail(ctx context.Context, email string) (string, error)
	RetrieveIdByUsername(ctx context.Context, username string) (string, error)
	UpdatePost(ctx context.Context, arg UpdatePostParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
