package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yudanl96/revive/util"
)

func CreateRandomPost(t *testing.T) (CreatePostParams, CreateUserParams) {
	userArg := CreateRandomUser(t)
	require.NotEmpty(t, userArg)

	id := uuid.New().String()
	postArg := CreatePostParams{
		ID:          id,
		UserID:      userArg.ID,
		Description: util.RandomLongStr(),
		Price:       int32(util.RandomInt(0, 200)),
	}

	err := testQueries.CreatePost(context.Background(), postArg)
	require.NoError(t, err)
	return postArg, userArg
}

func TestCreatePost(t *testing.T) {
	CreateRandomPost(t)
}

func TestGetPostById(t *testing.T) {
	postArg, _ := CreateRandomPost(t)

	post, err := testQueries.GetPostById(context.Background(), postArg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	require.Equal(t, postArg.ID, post.ID)
	require.Equal(t, postArg.Description, post.Description)
	require.Equal(t, postArg.Price, post.Price)
	require.Equal(t, postArg.UserID, post.UserID)

}

func TestUpdatePost(t *testing.T) {
	postArg, _ := CreateRandomPost(t)

	postArgNew := UpdatePostParams{
		ID: postArg.ID,
		Description: sql.NullString{
			String: util.RandomShortStr(),
			Valid:  true,
		},
		Sold: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Price: sql.NullInt32{
			Int32: int32(util.RandomInt(0, 200)),
			Valid: true,
		},
	}

	time.Sleep(10 * time.Second)
	err := testQueries.UpdatePost(context.Background(), postArgNew)
	require.NoError(t, err)
	postNew, err := testQueries.GetPostById(context.Background(), postArgNew.ID)
	require.NoError(t, err)
	require.Equal(t, postNew.Sold, postArgNew.Sold)
	require.Equal(t, postNew.Price, postArgNew.Price)
	require.Equal(t, postNew.Description, postArgNew.Description)

}

func TestDeletePost(t *testing.T) {
	postArg, _ := CreateRandomPost(t)

	err := testQueries.DeletePost(context.Background(), postArg.ID)
	require.NoError(t, err)

	post, err := testQueries.GetUserById(context.Background(), postArg.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, post)
}

func TestListPosts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomPost(t)
	}

	arg := ListPostsParams{
		Limit:  3,
		Offset: 7,
	}

	posts, err := testQueries.ListPosts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, posts, 3)

	for _, post := range posts {
		require.NotEmpty(t, post)
	}
}
