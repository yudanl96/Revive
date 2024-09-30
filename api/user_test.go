package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	mockDB "github.com/yudanl96/revive/db/mock"
	db "github.com/yudanl96/revive/db/sqlc"
	"github.com/yudanl96/revive/util"
)

func genRandomUser() db.User {
	return db.User{
		ID:       uuid.NewString(),
		Email:    util.RandomShortStr() + "@gmail.com",
		Username: util.RandomShortStr(),
		Password: util.HashPassword(util.RandomShortStr()),
	}
}

func TestGetUserByUsernameAPI(t *testing.T) {
	user := genRandomUser()

	testCases := []struct {
		name          string
		username      string
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user.ID, nil)

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:     "UsernameNotFound",
			username: user.Username,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return("", sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "IDNotFound",
			username: user.Username,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user.ID, nil)

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "FindIDInternalError",
			username: user.Username,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user.ID, nil)

				store.EXPECT().
					GetUserById(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "FindUserInternalError",
			username: user.Username,
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return("", sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidUsername",
			username: "1",
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					RetrieveIdByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		test_cur := testCases[i]

		t.Run(test_cur.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			test_cur.buildStubs(store)

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%v", test_cur.username)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			test_cur.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, gotUser, user)

}
