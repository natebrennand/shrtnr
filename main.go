package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"strings"
)

const (
	GET  string = "GET"
	POST string = "POST"
)

// route handler
type serverHandler struct {
	pool *redis.Pool
}

type statHandler serverHandler

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

func (s statHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != GET {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	shawty := apiHandler{s.pool.Get(), resp, strings.Replace(req.URL.Path, "/stats/", "", -1), ServerRequest{}}
	defer shawty.conn.Close()

	shawty.getUrlStats(resp, req)
}

// Routes all http requests to their corresponding functions
func (s serverHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	requestBody, err := getReqBody(req)
	// return error if POST request w/ no data
	if err != nil && req.Method == POST {
		resp.WriteHeader(http.StatusBadRequest)
		log.Printf("%d: Error decoding POST request's payload\n", http.StatusBadRequest)
		return
	}

	shawty := apiHandler{s.pool.Get(), resp, req.URL.Path[1:], requestBody}
	defer shawty.conn.Close()

	switch {
	// handle homepage & URL forwarding
	case req.Method == GET:
		switch url {
		case "/": // return homepage
			http.ServeFile(resp, req, "static/index.html")
		default:
			shawty.forward(*req, req.URL.Path[1:])
		}
	case req.Method == POST && url == "/":
		shawty.createShortUrl(resp, requestBody)
	default:
		// return a 501 error
		resp.WriteHeader(http.StatusNotImplemented)
		log.Printf("%d: Method not implemented, %s to %s", http.StatusNotImplemented, req.Method, url)
	}
}

func setupEndpoints() {
	// handles static asset packages & favicon
	http.Handle("/favion.ico", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, "static/favicon.ico")
	}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	rPool := redis.NewPool(redisPoolConnect, 2)

	// stat endpoint
	http.Handle("/stats/", statHandler{rPool})

	// handles all API routes
	http.Handle("/", serverHandler{rPool})
}

func main() {
	setupEndpoints()

	fmt.Println("running on 8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		panic("HTTP ListenAndServe failed")
		log.Fatal("HTTP ListenAndServe failed")
	}
}
