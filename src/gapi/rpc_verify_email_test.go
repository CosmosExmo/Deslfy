package gapi

import (
	"context"
	"database/sql"
	mockdb "desly/db/mock"
	db "desly/db/sqlc"
	"desly/pb"
	"desly/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomVerifyEmail(t *testing.T, owner string) db.VerifyEmail {
	return db.VerifyEmail{
		ID:         int64(util.RandomInt(10, 100)),
		Username:   owner,
		Email:      util.RandomEmail(),
		SecretCode: util.RandomString(32),
		IsUsed:     false,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(24 * time.Minute),
	}
}

func TestVerifyEmailGRPC(t *testing.T) {
	user, _ := randomUser(t)
	randomVerifyEmail := randomVerifyEmail(t, user.Username)

	testCases := []struct {
		name          string
		req           *pb.VerifyEmailRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.VerifyEmailRequest{
				EmailId:    randomVerifyEmail.ID,
				SecretCode: randomVerifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.VerifyEmailTxParams{
					EmailID:    randomVerifyEmail.ID,
					SecretCode: randomVerifyEmail.SecretCode,
				}

				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.OK, st.Code())
			},
		},
		{
			name: "ValidationError",
			req: &pb.VerifyEmailRequest{
				EmailId:    -1000,
				SecretCode: "123-22",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NotNil(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.VerifyEmailRequest{
				EmailId:    randomVerifyEmail.ID,
				SecretCode: randomVerifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.VerifyEmailTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
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

			res, err := server.VerifyEmail(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
