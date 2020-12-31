package auth

import (
	"math/rand"
)

const randomKeyTable string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@&*"
const tableLength int = len(randomKeyTable)

// generate random key
func GenerateKey(len int) string {
	buf := make([]rune, len)

	for i := 0; i < len; i++ {
		buf[i] = rune(randomKeyTable[rand.Intn(tableLength)])
	}

	return string(buf)
}
