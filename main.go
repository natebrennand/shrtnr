package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

const (
	GET  string = "GET"
	POST string = "POST"
)

// route handler
type serverHandler struct {
	pool *redis.Pool
}

// create a pool of redis connections
func redisPoolConnect() (redis.Conn, error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic("Cannot connect to Redis")
	}
	return c, err
}

type apiHandler struct {
	conn        redis.Conn
	resp        http.ResponseWriter
	shortURL    string
	requestBody ServerRequest
}

// Parses the JSON request body from the request
func getReqBody(req *http.Request) (ServerRequest, error) {
	decoder := json.NewDecoder(req.Body)
	var requestBody ServerRequest
	err := decoder.Decode(&requestBody)
	return requestBody, err
}

// Routes all http requests to their corresponding functions
func (s serverHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	shortURL := url[1:] // removes the initial '/'

	requestBody, err := getReqBody(req)
	// return error if POST request w/ no data
	if err != nil && req.Method == POST {
		resp.WriteHeader(http.StatusBadRequest)
		log.Printf("%d: Error decoding POST request's payload\n", http.StatusBadRequest)
		return
	}

	shawty := apiHandler{s.pool.Get(), resp, shortURL, requestBody}
	defer shawty.conn.Close()

	switch {
	case req.Method == GET:
		switch url {
		case "/": // return homepage
			http.ServeFile(resp, req, "static/index.html")
		case "/favicon.ico":
			http.ServeFile(resp, req, "static/favicon.ico")
			resp.WriteHeader(http.StatusNotImplemented)
		default:
			shawty.getLongUrl(*req, shortURL)
		}
	case req.Method == POST && url == "/":
		shawty.createShortUrl(resp, requestBody)
	default:
		// return a 501 error
		resp.WriteHeader(http.StatusNotImplemented)
		log.Printf("%d: Method not implemented, %s to %s", http.StatusNotImplemented, req.Method, url)
	}
}

func main() {
	// handles static asset packages
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// handles all API routes
	http.Handle("/", serverHandler{redis.NewPool(redisPoolConnect, 2)})

	fmt.Println("running on 8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		panic("HTTP ListenAndServe failed")
		log.Fatal("HTTP ListenAndServe failed")
	}
}
