package helpers

import (
	"hash/fnv"
	"math/rand"
	"time"
)

func HashString(s string) (uint64, error) {
	h := fnv.New64()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

func RandomToken(n int) string {
	characters := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
		"l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v",
		"w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6",
		"7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H",
		"I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S",
		"T", "U", "V", "W", "X", "Y", "Z",
	}
	var token string
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		token += characters[rand.Intn(len(characters))]
	}
	return token
}
