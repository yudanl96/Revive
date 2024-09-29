package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateSavedPost(t *testing.T) (CreatePostParams, CreateUserParams) {
	postArg, userArg := CreateRandomPost(t)

	arg := CreateSavedPostParams{
		UserID: userArg.ID,
		PostID: postArg.ID,
	}

	err := testQueries.CreateSavedPost(context.Background(), arg)
	require.NoError(t, err)
	return postArg, userArg
}

func TestCreateSavedPost(t *testing.T) {
	CreateSavedPost(t)
}

func TestGetSavedPostByIds(t *testing.T) {
	postArg, userArg := CreateSavedPost(t)

	arg := GetSavedPostByIdsParams{
		UserID: userArg.ID,
		PostID: postArg.ID,
	}

	savedPost, err := testQueries.GetSavedPostByIds(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, savedPost.PostID, postArg.ID)
	require.Equal(t, savedPost.UserID, userArg.ID)
}

func TestDeleteSavedPost(t *testing.T) {
	postArg, userArg := CreateSavedPost(t)

	arg := DeleteSavedPostParams{
		UserID: userArg.ID,
		PostID: postArg.ID,
	}

	err := testQueries.DeleteSavedPost(context.Background(), arg)
	require.NoError(t, err)

	argGet := GetSavedPostByIdsParams{
		UserID: userArg.ID,
		PostID: postArg.ID,
	}
	savedPost, err := testQueries.GetSavedPostByIds(context.Background(), argGet)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, savedPost)
}

func TestCountSavedPostsByPost(t *testing.T) {
	postArg, _ := CreateSavedPost(t)

	count, err := testQueries.CountSavedPostsByPost(context.Background(), postArg.ID)
	require.NoError(t, err)
	require.Equal(t, int64(count), int64(1))
}

func TestListSavedPostsByUser(t *testing.T) {
	_, userArg := CreateSavedPost(t)

	savedPosts, err := testQueries.ListSavedPostsByUser(context.Background(), userArg.ID)
	require.NoError(t, err)
	require.Len(t, savedPosts, 1)
	require.NotEmpty(t, savedPosts[0])
}
