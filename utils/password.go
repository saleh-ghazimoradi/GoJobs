package utils

import (
	"math/rand"
	"strings"
)

func GeneratePassword(d int) string {
	const digit = "0123456789abcdef"

	var password strings.Builder
	password.Grow(d)
	for i := 0; i < d; i++ {
		password.WriteByte(digit[rand.Intn(len(digit))])
	}
	return password.String()
}
