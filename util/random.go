package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	ownerLength    = 6
	currenciesPool = []string{"EUR", "USD", "CAD", "AUD", "GBP", "BTC", "NIG"}
	pool           = "aqwertyuioplkjhgfdsazxcvbm"
	minBalance     = 10
	maxBalance     = 1000
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(pool)
	for i := 0; i < n; i++ {
		index := rand.Intn(k)
		sb.WriteByte(pool[index])
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(ownerLength)
}

func RandomCurency() string {
	n := len(currenciesPool)
	index := rand.Intn(n)

	return currenciesPool[index]
}

func RandomBalance() int64 {
	return RandomInt(int64(minBalance), int64(maxBalance))
}
