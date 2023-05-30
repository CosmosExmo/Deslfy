package db

import (
	"context"
	"desly/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)

	desly, err := testQueries.CreateDesly(context.Background(), randomRedirect)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}

/* func TestGetDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), randomRedirect)

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	desly, err := testQueries.GetDesly(context.Background(), createdDesly.ID)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
} */

func TestGetDesly(t *testing.T) {
	var randomRedirect = util.RandomString(10)

	createdDesly, errCreate := testQueries.CreateDesly(context.Background(), randomRedirect)

	require.NoError(t, errCreate)
	require.NotEmpty(t, createdDesly)

	desly, err := testQueries.GetDesly(context.Background(), createdDesly.Desly)

	require.NoError(t, err)
	require.NotEmpty(t, desly)

	require.NotZero(t, desly.ID)
	require.NotZero(t, desly.CreatedAt)

	require.NotEmpty(t, desly.Redirect)
	require.NotEmpty(t, desly.Desly)
}