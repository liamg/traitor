package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var filenameRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890_-")

func Filename() string {
	b := make([]rune, 8+rand.Intn(8))
	for i := range b {
		b[i] = filenameRunes[rand.Intn(len(filenameRunes))]
	}
	return string(b)
}

var imageRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func Image() string {
	b := make([]rune, 8+rand.Intn(8))
	for i := range b {
		b[i] = imageRunes[rand.Intn(len(imageRunes))]
	}
	return string(b)
}
