package gapi

import (
	"context"
	"database/sql"
	mockdb "desly/db/mock"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/token"
	"desly/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomRefreshToken(t *testing.T, owner string) string {
	maker, err := token.NewPasetoMaker(testTokenSymmetricKey)
	require.NoError(t, err)

	token, _, err := maker.CreateToken(owner, time.Minute, token.RefreshToken)
	require.NoError(t, err)

	return token
}

func randomSession(t *testing.T, owner string) db.Session {
	return db.Session{
		ID:           uuid.New(),
		Username:     owner,
		RefreshToken: randomRefreshToken(t, owner),
		UserAgent:    util.RandomString(12),
		ClientIp:     util.RandomString(12),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(24 * time.Minute),
		CreatedAt:    time.Now(),
	}
}

func TestRenewAccessGRPC(t *testing.T) {
	user, _ := randomUser(t)
	randomSession := randomSession(t, user.Username)

	testCases := []struct {
		name          string
		req           *pb.RenewAccessRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.RenewAccessResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.RenewAccessRequest{
				RefreshToken: randomSession.RefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(randomSession, nil)
			},
			checkResponse: func(t *testing.T, res *pb.RenewAccessResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdAccessToken := res.GetAccessToken()
				require.NotEmpty(t, createdAccessToken)
			},
		},
		{
			name: "ValidationError",
			req: &pb.RenewAccessRequest{
				RefreshToken: "_???",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.RenewAccessResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.RenewAccessRequest{
				RefreshToken: randomSession.RefreshToken,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Session{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.RenewAccessResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			res, err := server.RenewAccess(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
