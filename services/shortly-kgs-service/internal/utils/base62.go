package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomKey(length int) (string, error) {

	var result strings.Builder

	base := int64(len(charset)) // 62 characters (Base62)

	// Generate a random number in a large range
	randNum, err := rand.Int(rand.Reader, big.NewInt(base))

	if err != nil {
		return "", err
	}

	for randNum.Cmp(big.NewInt(0)) == 1 {
		remainder := randNum.Int64() % base
		result.WriteByte(charset[remainder])
		randNum.Div(randNum, big.NewInt(base))
	}

	key := result.String()
	for len(key) < length {
		key = string(charset[0]) + key
	}

	return key, nil
}
