package api

import (
	"bytes"
	"database/sql"
	mockdb "desly/db/mock"
	db "desly/db/sqlc"
	"desly/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateDesly(t *testing.T) {
	desly := randomDesly()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"redirect": desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Eq(desly.Redirect)).
					Times(1).Return(desly, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchDesly(t, recorder.Body, desly)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"redirect": "0",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"redirect": desly.Redirect,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDesly(gomock.Any(), gomock.Eq(desly.Redirect)).
					Times(1).Return(db.Desly{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/desly"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetDesly(t *testing.T) {
	desly := randomDesly()

	testCases := []struct {
		name          string
		desly         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Eq(desly.Desly)).
					Times(1).Return(desly, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchDesly(t, recorder.Body, desly)
			},
		},
		{
			name:  "NotFound",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Eq(desly.Desly)).
					Times(1).Return(db.Desly{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "BadRequest",
			desly: "0",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "InternalServerError",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), desly.Desly).
					Times(1).Return(db.Desly{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/desly/%s", tc.desly)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestRedirect(t *testing.T) {
	desly := randomDesly()

	testCases := []struct {
		name          string
		desly         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Eq(desly.Desly)).
					Times(1).Return(desly, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Eq(desly.Desly)).
					Times(1).Return(db.Desly{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "BadRequest",
			desly: "0",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "InternalServerError",
			desly: desly.Desly,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetDesly(gomock.Any(), desly.Desly).
					Times(1).Return(db.Desly{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/r/%s", tc.desly)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomDesly() db.Desly {
	return db.Desly{
		Redirect: util.RandomString(20),
		Desly:    util.RandomString(6),
	}
}

func requireBodyMatchDesly(t *testing.T, body *bytes.Buffer, desly db.Desly) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotDesly db.Desly
	err = json.Unmarshal(data, &gotDesly)
	require.NoError(t, err)
	require.Equal(t, desly.Redirect, gotDesly.Redirect)
}
