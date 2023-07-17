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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func randomUserToken(owner string) db.UserToken {
	return db.UserToken{
		ID:        util.RandomInt(10, 50),
		Owner:     owner,
		Token:     util.RandomString(64),
		ExpireAt:  time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

func TestCreateUserTokenGRPC(t *testing.T) {
	user, _ := randomUser(t)
	randomToken := randomUserToken(user.Username)
	expireAt := timestamppb.New(randomToken.ExpireAt)

	invalidExpireAt := timestamppb.New(time.Now().Add(- (24*time.Hour)))

	testCases := []struct {
		name          string
		req           *pb.CreateUserTokenRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateUserTokenResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserTokenRequest{
				ExpireAt: expireAt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserToken(gomock.Any(), gomock.Any()).
					Times(1).Return(randomToken, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUserToken := res.GetUserToken()
				require.Equal(t, user.Username, createdUserToken.Owner)
				require.Equal(t, randomToken.Token, createdUserToken.Token)
				require.WithinDuration(t, randomToken.ExpireAt, randomToken.ExpireAt, time.Second)
			},
		},
		{
			name: "ValidationError",
			req: &pb.CreateUserTokenRequest{
				ExpireAt: invalidExpireAt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.CreateUserTokenRequest{
				ExpireAt: expireAt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.CreateUserTokenRequest{
				ExpireAt: expireAt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserToken(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserTokenResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserTokenRequest{
				ExpireAt: expireAt,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUserToken(gomock.Any(), gomock.Any()).
					Times(1).Return(db.UserToken{}, sql.ErrConnDone)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserTokenResponse, err error) {
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
			res, err := server.CreateUserToken(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
