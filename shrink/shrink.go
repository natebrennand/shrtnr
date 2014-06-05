package shrink

import (
	"errors"
	"math/rand"

	"github.com/garyburd/redigo/redis"
)

const (
	ALPHABET    string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	HASH_LENGTH int    = 5
	LONG        string = "LongURL"
	COUNT       string = "HitCount"
)

var (
	UrlInUse    error = errors.New("Short URL already in use")
	UrlNotFound error = errors.New("URL not found")
)

// returns a randomly generated shortened URL
func randURL() string {
	urlHash := ""
	for i := 0; i < HASH_LENGTH; i++ {
		randomChar := ALPHABET[int(rand.Float32()*float32(len(ALPHABET)))]
		urlHash += string(randomChar)
	}
	return urlHash
}

// creates the requested shortened URL
func CreateURL(conn redis.Conn, longURL string, shortURL string) (string, error) {
	// randomly assign URL if not given
	if shortURL == "" {
		for { // loop until unique string
			shortURL = randURL()
			v, err := redis.Int(conn.Do("EXISTS", shortURL))
			if err == nil && v == 0 {
				break
			}
		}
	} else { // confirm that URL is free
		v, err := redis.Int(conn.Do("EXISTS", shortURL))
		if err != nil {
			return "", err
		} else if v == 1 {
			return "", UrlInUse
		}
	}
	// create hash record
	v, err := redis.String(conn.Do("HMSET", shortURL, LONG, longURL, COUNT, 0))
	if v != "OK" || err != nil {
		return "", err
	}

	return shortURL, err
}

// retrieves a URL based on the shortened URL
func RetrieveURL(conn redis.Conn, shortURL string) (string, error) {
	// lookup long URL
	longURL, err := redis.String(conn.Do("HGET", shortURL, LONG))
	if err != nil {
		return "", UrlNotFound
	}

	// increment counter
	currentCount, err := redis.Int(conn.Do("HINCRBY", shortURL, COUNT, 1))
	if err != nil || currentCount < 0 {
		return "", UrlNotFound
	}
	return longURL, nil
}
