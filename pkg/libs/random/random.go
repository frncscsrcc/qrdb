package random

import (
	"math/rand"
	"time"
)

func GetRandomString(lenght int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, lenght)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
