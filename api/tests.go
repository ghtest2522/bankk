package api

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/util"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase[T any] struct {
	name          string
	data          T
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, message string)
}

func CreateRandomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurency(),
	}
}

func GetCreateAccountReq() CreateAccountRequest {
	return CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurency(),
	}
}

func GetCreateAccountInvlidOwnerReq() CreateAccountRequest {
	return CreateAccountRequest{
		Owner:    "",
		Currency: util.RandomCurency(),
	}
}

func GetCreateAccountInvlidCurrencyReq() CreateAccountRequest {
	return CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: "INVALID",
	}
}

func GetCreateAccountMissingCurrencyReq() CreateAccountRequest {
	return CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: "",
	}
}

func ReadBodyRespone[T any](t *testing.T, body *bytes.Buffer) ResponseResult[T] {
	reader, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var result ResponseResult[T]
	err = json.Unmarshal(reader, &result)
	require.NoError(t, err)

	return result
}
