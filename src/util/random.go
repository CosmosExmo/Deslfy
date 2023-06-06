package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//RandomString generates a random string of length n
func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//Random Int generates a random int between min and max
func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min + 1)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}