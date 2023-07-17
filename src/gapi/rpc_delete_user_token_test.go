package gapi

import (
	"context"
	"database/sql"
	mockdb "desly/db/mock"
	"desly/pb"
	"desly/token"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func TestDeleteUserTokenGRPC(t *testing.T) {
	user, _ := randomUser(t)
	randomToken := randomUserToken(user.Username)

	testCases := []struct {
		name          string
		req           *pb.DeleteUserTokenRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.DeleteUserTokenResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteUserTokenRequest{
				Id: randomToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserToken(gomock.Any(), gomock.Any()).
					Times(1).Return(nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				deleteResponse := res.GetIsDeleteSuccessful()
				require.True(t, deleteResponse)
			},
		},
		{
			name: "ValidationError",
			req: &pb.DeleteUserTokenRequest{
				Id: -1000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.DeleteUserTokenRequest{
				Id: randomToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.DeleteUserTokenRequest{
				Id: randomToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.DeleteUserTokenRequest{
				Id: randomToken.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUserToken(gomock.Any(), gomock.Any()).
					Times(1).Return(sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserTokenResponse, err error) {
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
			res, err := server.DeleteUserToken(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}