package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/util"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type GetAccountData struct {
	accountID int64
	message   string
}

func TestGetAccountAPI(t *testing.T) {
	account := CreateRandomAccount()

	testCases := []TestCase[GetAccountData]{
		{
			name: "OK",
			data: GetAccountData{message: util.AccountWasFound, accountID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusOK)

				result := ReadBodyRespone[db.Account](t, recorder.Body)
				require.Equal(t, account, result.Data)
				require.Equal(t, util.AccountWasFound, result.Message)
			},
		},
		{
			name: "Bad Account ID",
			data: GetAccountData{accountID: 0},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusBadRequest)

				result := ReadBodyRespone[any](t, recorder.Body)
				require.Empty(t, result.Data)
				require.NotZero(t, result.Message)
			},
		},
		{
			name: "Not Found Account",
			data: GetAccountData{accountID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusNotFound)

				result := ReadBodyRespone[any](t, recorder.Body)
				require.Empty(t, result.Data)
				require.NotZero(t, result.Message)
				require.Equal(t, sql.ErrNoRows.Error(), result.Message)
			},
		},
		{
			name: "Internal Server Error",
			data: GetAccountData{accountID: account.ID},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusInternalServerError)

				result := ReadBodyRespone[any](t, recorder.Body)
				require.Empty(t, result.Data)
				require.NotZero(t, result.Message)
				require.Equal(t, sql.ErrConnDone.Error(), result.Message)
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
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// make request
			url := fmt.Sprintf("/account/%d", tc.data.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			// check response
			tc.checkResponse(t, recorder, tc.data.message)
		})
	}

}
