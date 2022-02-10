package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name          string
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, message string)
	accountID     int64
	message       string
}

func TestGetAccountAPI(t *testing.T) {
	account := createRandomAccount()

	testCases := []TestCase{
		{
			name:      "OK",
			message:   util.AccountWasFound,
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusOK)
				checkForBodyMatch(t, recorder.Body, account, message)
			},
		},
		{
			name:      "Bad Account ID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(gomock.Any())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusBadRequest)
			},
		},
		{
			name:      "Not Found Account",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusNotFound)
			},
		},
		{
			name:      "Internal Server Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusInternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// prepare server
			server := NewServer(store, util.NewJSONResponseSender())
			recorder := httptest.NewRecorder()

			// make request
			url := fmt.Sprintf("/account/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			// check response
			tc.checkResponse(t, recorder, tc.message)
		})
	}

}

func checkForBodyMatch(t *testing.T, body *bytes.Buffer, accout db.Account, message string) {
	reader, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var result db.Account
	err = json.Unmarshal(reader, &result)
	require.NoError(t, err)
	require.Equal(t, accout, result)
}

func createRandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurency(),
	}
}
