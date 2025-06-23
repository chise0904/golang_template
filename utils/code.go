package util

import (
	"math/rand"
	"time"
)

func GenRandomCode(length int) string {

	codeVal := ""
	codeRange := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	source := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(source)
	size := len(codeRange)
	for i := 0; i < length; i++ {
		codeVal += string(codeRange[newRand.Intn(size)])
	}

	return codeVal
}
