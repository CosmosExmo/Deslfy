package db

import (
	"context"
	"desly/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)
	var randomDesly = util.RandomString(6)

	desly, err := testQueries.CreateDesly(context.Background(), CreateDeslyParams{
		Redirect: randomRedirect,
		Desly: randomDesly,
	})

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}

func TestGetDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)
	var randomDesly = util.RandomString(6)

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), CreateDeslyParams{
		Redirect: randomRedirect,
		Desly: randomDesly,
	})

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	desly, err := testQueries.GetDesly(context.Background(), createdDesly.ID)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}

func TestGetRedirectByDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)
	var randomDesly = util.RandomString(6)

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), CreateDeslyParams{
		Redirect: randomRedirect,
		Desly: randomDesly,
	})

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	desly, err := testQueries.GetRedirectByDesly(context.Background(), createdDesly.Desly)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}