package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateRandomNumber generates a random integer between min and max (inclusive).
func GenerateRandomNumber(min, max int) (int, error) {
	if min > max {
		return 0, fmt.Errorf("min must be less than or equal to max")
	}

	num, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}

	return int(num.Int64()) + min, nil
}
