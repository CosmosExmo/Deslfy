package gapi

import (
	db "desly/db/sqlc"
	"desly/util"
	"desly/worker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}