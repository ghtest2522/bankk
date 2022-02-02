package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)
	fmt.Println("Balances before: ", acc1.Balance, acc2.Balance)

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for j := 0; j < n; j++ {
		err := <-errs
		result := <-results

		// No Error occourd
		require.NoError(t, err)
		// Results are not empty
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.NotEmpty(t, result.FromEntry)
		require.NotEmpty(t, result.ToEntry)
		// Response is expected
		require.Equal(t, result.Transfer.FromAccountID, acc1.ID)
		require.Equal(t, result.Transfer.ToAccountID, acc2.ID)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)
		// Data is presistent: Transfer
		transfer, err := store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
		// Data is presistent: FromEntry
		entry, err := store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)
		require.Equal(t, entry.ID, result.FromEntry.ID)
		require.Equal(t, entry.Amount, result.FromEntry.Amount)
		require.Equal(t, entry.CreatedAt, result.FromEntry.CreatedAt)
		require.Equal(t, entry.AccountID, result.FromEntry.AccountID)
		// Data is presistent: ToEntry
		entry, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		require.Equal(t, entry.ID, result.ToEntry.ID)
		require.Equal(t, entry.Amount, result.ToEntry.Amount)
		require.Equal(t, entry.CreatedAt, result.ToEntry.CreatedAt)
		require.Equal(t, entry.AccountID, result.ToEntry.AccountID)
		//
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acc1.ID, fromAccount.ID)
		//
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acc2.ID, toAccount.ID)
		//
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		fmt.Println("Balances tx: ", fromAccount.Balance, toAccount.Balance)
		fmt.Println("Diff tx: ", diff1, diff2)
		require.Equal(t, diff2, diff1)
		require.Zero(t, diff1%amount)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && n >= k)

	}
	//
	acc1Update, err1 := store.GetAccount(context.Background(), acc1.ID)
	acc2Update, err := store.GetAccount(context.Background(), acc2.ID)
	fmt.Println("Balances after: ", acc1Update.Balance, acc2Update.Balance)
	require.NoError(t, err1)
	expectedBalance := acc1.Balance - int64(n*int(amount))
	require.Equal(t, acc1Update.Balance, expectedBalance)
	//
	require.NoError(t, err)
	expectedBalance = acc2.Balance + int64(n*int(amount))
	require.Equal(t, acc2Update.Balance, expectedBalance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	acc1 := CreateRandomAccount(t)
	acc2 := CreateRandomAccount(t)

	n := 18
	amount := int64(10)
	errs := make(chan error)

	fmt.Println("Acc1: ", acc1.Balance, " Acc2: ", acc2.Balance)

	for i := 0; i < n; i++ {
		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 0 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}
		fmt.Println("Transfering ", amount, " From: ", fromAccountID, " To: ", toAccountID)
		go func() {

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for j := 0; j < n; j++ {
		err := <-errs
		require.NoError(t, err)
	}
	//

	acc1Update, _ := store.GetAccount(context.Background(), acc1.ID)
	acc2Update, _ := store.GetAccount(context.Background(), acc2.ID)
	fmt.Println("Acc1: ", acc1Update.Balance, " Acc2: ", acc2Update.Balance)

	require.Equal(t, acc1Update.Balance, acc1.Balance)
	require.Equal(t, acc2.Balance, acc2Update.Balance)

}
