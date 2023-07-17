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
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomDesly(owner string) db.Desly {
	return db.Desly{
		Redirect: util.RandomURL(),
		Desly:    util.RandomString(6),
		Owner:    owner,
	}
}

func TestCreateDeslyGRPC(t *testing.T) {
	user, _ := randomUser(t)
	desly := randomDesly(user.Username)

	invalidRedirect := util.RandomString(5)

	testCases := []struct {
		name          string
		req           *pb.CreateDeslyRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateDeslyResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateDeslyRequest{
				Redirect: desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateDeslyParams{
					Redirect: desly.Redirect,
					Owner:    user.Username,
				}

				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(desly, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateDeslyResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdDesly := res.GetDesly()
				require.Equal(t, user.Username, createdDesly.Owner)
				require.Equal(t, desly.Redirect, createdDesly.Redirect)
				require.Equal(t, desly.Desly, createdDesly.Desly)
			},
		},
		{
			name: "ValidationError",
			req: &pb.CreateDeslyRequest{
				Redirect: invalidRedirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateDeslyResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.CreateDeslyRequest{
				Redirect: desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateDeslyResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.CreateDeslyRequest{
				Redirect: desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.CreateDeslyResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateDeslyRequest{
				Redirect: desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Desly{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateDeslyResponse, err error) {
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
			res, err := server.CreateDesly(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
