package shrink

import (
	"errors"
	"github.com/garyburd/redigo/redis"
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
)

// returns a new redis connection
func RedisConn() redis.Conn {
	// TODO: pool connections in some way, don't connect each request
	conn, err := redis.Dial(NETWORK, ":"+strconv.Itoa(REDIS_PORT))
	if err != nil {
		panic(err.Error())
	}
	return conn
}

// returns a randomly generated shortened URL
func RandURL() string {
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
	conn := RedisConn()
	defer conn.Close()
	if shortURL == "" {
		shortURL = RandURL()
	} else {
		v, err := redis.Int(conn.Do("EXISTS", shortURL))
		if err != nil {
			return "", err
		} else if v == 1 {
			return "", UrlInUse
		}
	}
	v, err := conn.Do("SET", shortURL, longURL)
	if v != "OK" || err != nil {
		return "", err
	}
	return shortURL, err
}

// retrieves a URL based on the shortened URL
func RetrieveURL(shortURL string) (string, error) {
	conn := RedisConn()
	defer conn.Close()
	longURL, err := redis.String(conn.Do("GET", shortURL))
	if err != nil {
		return "", UrlNotFound
	}
	return longURL, nil
}
