package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type CreateAccountData struct {
	payload CreateAccountRequest
	message string
}

func TestCreateAccountAPI(t *testing.T) {
	createAccountGoodPayload := GetCreateAccountReq()
	okResult := db.Account{Owner: createAccountGoodPayload.Owner, Currency: createAccountGoodPayload.Currency, ID: util.RandomInt(1, 1000), Balance: 0}

	testCases := []TestCase[CreateAccountData]{
		{
			name: "OK",
			data: CreateAccountData{payload: createAccountGoodPayload, message: util.AccountWasCeated},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(db.CreateAccountParams{
						Owner:    createAccountGoodPayload.Owner,
						Currency: createAccountGoodPayload.Currency,
						Balance:  0,
					})).
					Times(1).Return(okResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusOK)

				result := ReadBodyRespone[db.Account](t, recorder.Body)
				require.Equal(t, okResult, result.Data)
				require.Equal(t, util.AccountWasCeated, result.Message)
			},
		},
		{
			name: "Invalid Account Owner",
			data: CreateAccountData{payload: GetCreateAccountInvlidOwnerReq()},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(gomock.Any())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusBadRequest)
				result := ReadBodyRespone[any](t, recorder.Body)
				require.NotEmpty(t, result.Message)
				require.Empty(t, result.Data)
			},
		},
		{
			name: "Missing Currency",
			data: CreateAccountData{payload: GetCreateAccountMissingCurrencyReq()},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(gomock.Any())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				result := ReadBodyRespone[any](t, recorder.Body)
				require.NotEmpty(t, result.Message)
				require.Empty(t, result.Data)
			},
		},
		{
			name: "Invalid Currency",
			data: CreateAccountData{payload: GetCreateAccountInvlidCurrencyReq()},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(gomock.Any())).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
				result := ReadBodyRespone[any](t, recorder.Body)
				require.NotEmpty(t, result.Message)
				require.Empty(t, result.Data)
			},
		},
		{
			name: "Internal Server Error",
			data: CreateAccountData{payload: createAccountGoodPayload},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
				statusCode := recorder.Code
				require.Equal(t, statusCode, http.StatusInternalServerError)

				result := ReadBodyRespone[any](t, recorder.Body)
				require.Equal(t, sql.ErrConnDone.Error(), result.Message)
				require.Empty(t, result.Data)
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
			url := fmt.Sprintf("/accounts")
			payload, err := json.Marshal(tc.data.payload)
			body := bytes.NewReader(payload)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, req)

			// check response
			tc.checkResponse(t, recorder, tc.data.message)
		})
	}

}
