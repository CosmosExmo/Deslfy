package gapi

import (
	"context"
	"database/sql"
	mockdb "desly/db/mock"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/token"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUserTokensGRPC(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	userTokens := make([]db.UserToken, n)
	for i := 0; i < n; i++ {
		userTokens[i] = randomUserToken(user.Username)
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.GetUserTokensResponse, err error)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetUserTokens(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).Return(userTokens, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokensResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUserTokens := res.GetUserTokens()
				for _, userToken := range createdUserTokens {
					require.Equal(t, user.Username, userToken.Owner)
					require.NotEmpty(t, userToken.Token)
					require.NotEmpty(t, userToken.ExpireAt)
					require.NotEmpty(t, userToken.CreatedAt)
				}
			},
		},
		{
			name: "ExpiredToken",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserTokens(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokensResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserTokens(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokensResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserTokens(gomock.Any(), gomock.Any()).
					Times(1).Return([]db.UserToken{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokensResponse, err error) {
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

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.GetUserTokens(ctx, &pb.GetUserTokensRequest{})
			tc.checkResponse(t, res, err)
		})
	}
}
