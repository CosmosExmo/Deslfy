package db

import (
	"context"
	"desly/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDesly(t *testing.T) {
	account := createRandomUser(t)
	var randomRedirect = util.RandomString(10)

	arg := CreateDeslyParams{
		Redirect: randomRedirect,
		Owner:    account.Username,
	}

	desly, err := testQueries.CreateDesly(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}

func TestGetDesly(t *testing.T) {
	account := createRandomUser(t)
	var randomRedirect = util.RandomString(10)

	arg := CreateDeslyParams{
		Redirect: randomRedirect,
		Owner:    account.Username,
	}

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), arg)

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	getArg := GetDeslyParams{
		Desly: createdDesly.Desly,
		Owner: account.Username,
	}
	desly, err := testQueries.GetDesly(context.Background(), getArg)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}

func TestGetRedirectByDesly(t *testing.T) {
	account := createRandomUser(t)
	var randomRedirect = util.RandomString(10)

	arg := CreateDeslyParams{
		Redirect: randomRedirect,
		Owner:    account.Username,
	}

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), arg)

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	redirect, err := testQueries.GetRedirectByDesly(context.Background(), createdDesly.Desly)

	require.NoError(t, err)
	require.NotEmpty(t, redirect)
}
