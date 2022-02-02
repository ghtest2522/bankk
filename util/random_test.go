package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	const (
		max = 1000
		min = 1
	)

	result := RandomInt(min, max)
	require.Less(t, result, int64(max))
	require.GreaterOrEqual(t, result, int64(min))
}

func TestRandomString(t *testing.T) {
	const (
		len = 5
	)
	result := RandomString(len)
	require.Len(t, result, len)

	for _, char := range result {
		require.Contains(t, pool, string(char))
	}
}

func TestRandomOwner(t *testing.T) {
	result := RandomOwner()
	require.Len(t, result, ownerLength)
}

func TestRandomCurency(t *testing.T) {
	result := RandomCurency()
	var resultBelongsToPool bool

	for _, v := range currenciesPool {
		if v == result {
			resultBelongsToPool = true
		}
	}

	require.True(t, resultBelongsToPool)
}

func TestRandomBalance(t *testing.T) {
	result := RandomBalance()
	require.Less(t, result, int64(maxBalance))
	require.GreaterOrEqual(t, result, int64(minBalance))
}
