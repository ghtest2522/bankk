package db

import (
	"bank/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurency(),
	}

	account, err := testQuries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := CreateRandomAccount(t)
	acc2, err := testQuries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := CreateRandomAccount(t)

	update := UpdateAccountParams{
		Balance: acc.Balance + util.RandomBalance(),
		ID:      acc.ID,
	}
	result, err := testQuries.UpdateAccount(context.Background(), update)

	// Account was updated correctly
	require.NoError(t, err)
	require.Equal(t, result.Balance, update.Balance)

	// Rest of feilds are not affected
	require.Equal(t, acc.Currency, result.Currency)
	require.Equal(t, acc.ID, result.ID)
	require.Equal(t, acc.Owner, result.Owner)
	require.WithinDuration(t, acc.CreatedAt, result.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc := CreateRandomAccount(t)

	err := testQuries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	acc2, err := testQuries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQuries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
