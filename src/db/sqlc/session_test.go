package db

import (
	"context"
	"desly/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {
	user := createRandomUser(t)

	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: util.RandomString(64),
		UserAgent:    util.RandomString(12),
		ClientIp:     util.RandomString(24),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.NotZero(t, session.CreatedAt)

	require.Equal(t, session.Username, arg.Username)
	require.Equal(t, session.RefreshToken, arg.RefreshToken)
	require.Equal(t, session.UserAgent, arg.UserAgent)
	require.Equal(t, session.IsBlocked, arg.IsBlocked)
	require.Equal(t, session.ClientIp, arg.ClientIp)
	require.WithinDuration(t, session.ExpiresAt, arg.ExpiresAt, time.Second)

	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	randomSession := createRandomSession(t)

	session, err := testQueries.GetSession(context.Background(), randomSession.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.NotZero(t, session.CreatedAt)

	require.NotEmpty(t, session.ID)
	require.NotEmpty(t, session.Username)
	require.NotEmpty(t, session.RefreshToken)
	require.NotEmpty(t, session.UserAgent)
	require.NotEmpty(t, session.ClientIp)
	require.False(t, session.IsBlocked)
	require.NotEmpty(t, session.ExpiresAt)
}
