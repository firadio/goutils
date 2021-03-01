package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandInt64(min int64, max int64) int64 {
	count := max - min + 1
	b := new(big.Int).SetInt64(int64(count))
	i, err := rand.Int(rand.Reader, b)
	if err != nil {
		fmt.Printf("Can't generate random value: %v, %v", i, err)
		return 0
	}
	return min + i.Int64()
}
