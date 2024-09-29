package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yudanl96/revive/util"
)

func createRandomUser(t *testing.T) CreateUserParams {
	password := util.RandomStrLen()
	username := util.RandomStrLen()
	email := util.RandomStrLen() + "@gmail.com"
	hashedPW := util.HashPassword(password)
	err := util.MatchPassword(hashedPW, password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: username,
		Email:    email,
		Password: string(hashedPW),
	}

	err = testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	return arg
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestRetrieveIdByEmail(t *testing.T) {
	arg := createRandomUser(t)
	id, err := testQueries.RetrieveIdByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestRetrieveIdByUsername(t *testing.T) {
	arg := createRandomUser(t)
	id, err := testQueries.RetrieveIdByUsername(context.Background(), arg.Username)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestGetUserById(t *testing.T) {
	arg := createRandomUser(t)
	id, err := testQueries.RetrieveIdByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	user, err := testQueries.GetUserById(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
}

func TestUpdateUser(t *testing.T) {
	arg := createRandomUser(t)
	id, err := testQueries.RetrieveIdByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	arg_update := UpdateUserParams{
		Username: util.RandomStrLen(),
		Email:    util.RandomStrLen(),
		Password: arg.Password,
		ID:       id,
	}
	err = testQueries.UpdateUser(context.Background(), arg_update)
	require.NoError(t, err)
	user, err := testQueries.GetUserById(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, arg_update.Email, user.Email)
	require.Equal(t, arg_update.Username, user.Username)
	require.Equal(t, arg_update.Password, user.Password)
}

func TestDeleteUser(t *testing.T) {
	arg := createRandomUser(t)
	id, err := testQueries.RetrieveIdByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = testQueries.DeleteUser(context.Background(), id)
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), id)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  3,
		Offset: 7,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 3)

	for _, account := range users {
		require.NotEmpty(t, account)
	}
}
