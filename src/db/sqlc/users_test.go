package db

import (
	"context"
	"desly/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomString(8),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(15),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.FullName, arg.FullName)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	arg := CreateUserParams{
		Username:       util.RandomString(8),
		HashedPassword: "secret",
		FullName:       util.RandomString(15),
		Email:          util.RandomEmail(),
	}

	createdUser, errCreate := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdUser)

	user, err := testQueries.GetUser(context.Background(), arg.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotEmpty(t, user.Username)
	require.NotEmpty(t, user.HashedPassword)
	require.NotEmpty(t, user.Email)
	require.NotEmpty(t, user.FullName)
}