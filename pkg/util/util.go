package util

import (
	"crypto/rand"
	"io"
	"log"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// GenPinCode - generate random pin code
func GenPinCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		log.Fatal(err)
		return ""
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b)
}
