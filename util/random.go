package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Return random int betwen min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Return random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	alphabetLen := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(alphabetLen)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}

	currenciesLen := len(currencies)

	return currencies[rand.Intn(currenciesLen)]
}
