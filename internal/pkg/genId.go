package pkg

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenID(length int) (string, error) {
	seed := "012345679"
	byteSlice := make([]byte, length)

	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(seed)))

		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to generate random num: %w", err)
		}

		byteSlice[i] = seed[num.Int64()]
	}

	return string(byteSlice), nil
}
