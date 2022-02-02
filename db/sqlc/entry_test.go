package db

import (
	"bank/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)

	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	entry, err := testQuries.CreateEntry(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, entry.Amount, args.Amount)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	newEntry := createRandomEntry(t)
	entry, err := testQuries.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.Equal(t, entry, newEntry)
}

func TestListEntries(t *testing.T) {
	var lastEntry Entry
	for i := 0; i < 10; i++ {
		lastEntry = createRandomEntry(t)
	}

	args := ListEntriesParams{
		Limit:     5,
		Offset:    0,
		AccountID: lastEntry.AccountID,
	}
	entries, err := testQuries.ListEntries(context.Background(), args)
	require.NoError(t, err)

	for _, e := range entries {
		require.NotEmpty(t, e)
		require.Equal(t, e.AccountID, lastEntry.AccountID)
	}
}
