package db

import (
	"context"
	"database/sql"
	"desly/token"
	"desly/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomToken(t *testing.T, user User) UserToken {
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	duration := time.Hour
	accessToken, payload, err := maker.CreateToken(user.Username, duration, token.AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, payload)

	expireAt := time.Now().Add(duration)
	arg := CreateUserTokenParams{
		Owner: user.Username,
		Token: accessToken,
		ExpireAt: expireAt,
	}
	token, err := testQueries.CreateUserToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, token.ID)
	require.Equal(t, token.Owner, user.Username)
	require.Equal(t, token.Token, accessToken)
	require.WithinDuration(t, token.ExpireAt, expireAt, time.Second)

	return token
} 

func TestCreateUserToken(t *testing.T) {
	user := createRandomUser(t)
	createRandomToken(t, user)
}

func TestGetUserToken(t *testing.T) {
	user := createRandomUser(t)
	randomToken := createRandomToken(t, user)

	arg := GetUserTokenParams{
		Owner: randomToken.Owner,
		ID: randomToken.ID,
	}
	token, err := testQueries.GetUserToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.Equal(t, token, randomToken)
}

func TestGetUserTokens(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomToken(t, user)
	}

	tokens, err := testQueries.GetUserTokens(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, tokens)

	require.Len(t, tokens, 5)
	for _, token := range tokens {
		require.NotEmpty(t, token)
		require.Equal(t, token.Owner, user.Username)
	}
}

func TestDeleteUserToken(t *testing.T) {
	user := createRandomUser(t)
	randomToken := createRandomToken(t, user)
	
	deleteArg := DeleteUserTokenParams{
		ID: randomToken.ID,
		Owner: randomToken.Owner,
	}
	err := testQueries.DeleteUserToken(context.Background(), deleteArg)
	require.NoError(t, err)

	createArg := GetUserTokenParams{
		Owner: randomToken.Owner,
		ID: randomToken.ID,
	}
	_, err = testQueries.GetUserToken(context.Background(), createArg)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
}