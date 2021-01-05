package auth

import (
	"math/rand"
)

const randomKeyTable string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@&*"
const tableLength int = len(randomKeyTable)

var keyRegister map[string]string = make(map[string]string)

func RegisterKey(keyName string, keyValue string) {
	keyRegister[keyName] = keyValue
}

func GetKey(keyName string) string {
	return keyRegister[keyName]
}

// generate random key
func GenerateKey(len int) string {
	buf := make([]rune, len)

	for i := 0; i < len; i++ {
		buf[i] = rune(randomKeyTable[rand.Intn(tableLength)])
	}

	return string(buf)
}
