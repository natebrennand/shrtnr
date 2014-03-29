package shrink

import (
	"github.com/garyburd/redigo/redis"

	"errors"
	"math/rand"
	"strconv"
)

const (
	ALPHABET    string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	HASH_LENGTH int    = 5
	REDIS_PORT  int    = 6379
	NETWORK     string = "tcp"
)

var (
	UrlInUse    error = errors.New("Short URL already in use")
	UrlNotFound error = errors.New("URL not found")
	RedisConn redis.Conn
)

// returns a new redis connection
func Connect() {
	conn, err := redis.Dial(NETWORK, ":"+strconv.Itoa(REDIS_PORT))
	if err != nil {
		panic(err.Error())
	}
	RedisConn = conn
}

// returns a randomly generated shortened URL
func RandURL() string {
	urlHash := ""
	for i := 0; i < HASH_LENGTH; i++ {
		randomChar := ALPHABET[int(rand.Float32()*float32(len(ALPHABET)))]
		urlHash += string(randomChar)
	}
	return urlHash
}

// creates the requested shortened URL
func CreateURL(longURL string, shortURL string) (string, error) {
	if shortURL == "" {	// randomly assign URL
		for { // loop until unique
			shortURL = RandURL()
			v, err := redis.Int(RedisConn.Do("EXISTS", shortURL))
			if err == nil && v == 0 {
				break
			}
		}
	} else { // check that URL is free
		v, err := redis.Int(RedisConn.Do("EXISTS", shortURL))
		if err != nil {
			return "", err
		} else if v == 1 {
			return "", UrlInUse
		}
	}
	v, err := RedisConn.Do("SET", shortURL, longURL)
	if v != "OK" || err != nil {
		return "", err
	}
	return shortURL, err
}

// retrieves a URL based on the shortened URL
func RetrieveURL(shortURL string) (string, error) {
	longURL, err := redis.String(RedisConn.Do("GET", shortURL))
	if err != nil {
		return "", UrlNotFound
	}
	return longURL, nil
}
