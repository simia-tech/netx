package netx

import (
	"crypto/rand"
	"fmt"
)

// RandomAddress returns a random address with the provided prefix.
func RandomAddress(prefix string) string {
	bytes := [8]byte{}
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%x", prefix, bytes[:])
}
