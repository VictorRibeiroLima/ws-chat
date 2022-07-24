package util

import (
	"math/rand"
	"time"
)

func RandomString() string {
	rand.Seed(time.Now().Unix())

	str := "abcdefghijklmnopqrstuvxwyzABCDEFGHIJKLMNOPQRSTUVXWYZ"

	shuff := []rune(str)

	rand.Shuffle(len(shuff), func(i, j int) {
		shuff[i], shuff[j] = shuff[j], shuff[i]
	})
	return string(shuff)
}
