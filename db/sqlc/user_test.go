package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/yudanl96/revive/util"
)

func CreateRandomUser(t *testing.T) CreateUserParams {
	id := uuid.New().String()
	password := util.RandomShortStr()
	username := util.RandomShortStr()
	email := util.RandomShortStr() + "@gmail.com"

	arg := CreateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}

	err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	return arg
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestRetrieveIdByEmail(t *testing.T) {
	arg := CreateRandomUser(t)
	id, err := testQueries.RetrieveIdByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Equal(t, arg.ID, id)
}

func TestRetrieveIdByUsername(t *testing.T) {
	arg := CreateRandomUser(t)
	id, err := testQueries.RetrieveIdByUsername(context.Background(), arg.Username)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Equal(t, arg.ID, id)
}

func TestGetUserById(t *testing.T) {
	arg := CreateRandomUser(t)
	user, err := testQueries.GetUserById(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.ID, user.ID)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
}

func TestUpdateUser(t *testing.T) {
	arg := CreateRandomUser(t)
	arg_update := UpdateUserParams{
		Username: util.RandomShortStr(),
		Email:    util.RandomShortStr(),
		Password: arg.Password,
		ID:       arg.ID,
	}
	err := testQueries.UpdateUser(context.Background(), arg_update)
	require.NoError(t, err)
	user, err := testQueries.GetUserById(context.Background(), arg.ID)
	require.NoError(t, err)
	require.Equal(t, arg_update.Email, user.Email)
	require.Equal(t, arg_update.Username, user.Username)
	require.Equal(t, arg_update.Password, user.Password)
}

func TestDeleteUser(t *testing.T) {
	arg := CreateRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), arg.ID)
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), arg.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  3,
		Offset: 7,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 3)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
