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

func TestGetUserTokenGRPC(t *testing.T) {
	user, _ := randomUser(t)
	randomUserToken := randomUserToken(user.Username)

	testCases := []struct {
		name          string
		req           *pb.GetUserTokenRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.GetUserTokenResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetUserTokenRequest{
				Id: randomUserToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetUserTokenParams{
					ID:    randomUserToken.ID,
					Owner: user.Username,
				}

				store.EXPECT().
					GetUserToken(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(randomUserToken, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUserToken := res.GetUserToken()
				require.Equal(t, user.Username, createdUserToken.Owner)
				require.Equal(t, randomUserToken.ID, createdUserToken.Id)
				require.WithinDuration(t, randomUserToken.CreatedAt, createdUserToken.CreatedAt.AsTime(), time.Second)
			},
		},
		{
			name: "ValidationError",
			req: &pb.GetUserTokenRequest{
				Id: -1000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.GetUserTokenRequest{
				Id: randomUserToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.GetUserTokenRequest{
				Id: randomUserToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.GetUserTokenRequest{
				Id: randomUserToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserToken(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserToken{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.GetUserTokenResponse, err error) {
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
			res, err := server.GetUserToken(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
