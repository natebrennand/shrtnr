package shrink

import (
	"math/rand"
)

const (
	ALPHABET string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	HASH_LENGTH int = 5
)

// returns a randomly generated shortened URL
func RandURL() (string) {
	urlHash := ""
	for i := 0; i < HASH_LENGTH; i++ {
		randomChar := ALPHABET[int(rand.Float32()*float32(len(ALPHABET)))]
		urlHash += string(randomChar)
	}
	// TODO: do check w/ redis that new URL is unused
	return urlHash
}

// creates the requested shortened URL
func CreateURL(longURL string, shortURL string) (string, error) {
	if longURL == "" {
		longURL = RandURL()
	}
	return longURL, nil
}

// retrieves a URL based on the shortened URL
func RetrieveURL(shortURL string) (string, error) {
	return "", nil
}
