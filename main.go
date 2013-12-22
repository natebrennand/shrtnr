package main

import (
	"fmt"
	"net/http"
	"shrtnr/shrink"
)

const (
	GET string = "GET"
	POST string = "POST"
	PUT string = "PUT"
)

func router(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	method := req.Method
	shortURL := url[1:]

	switch {
	case url == "/": // homepage or creating a random url
		switch method {
		case GET:
			// returns static files
		case POST:
			// creates a random endpoint
		}
	case url != "/": // querying for a URL or creating a specific one
		switch method {
		case GET:
			// return the long url
			longURL, err := shrink.RetrieveURL(shortURL)
			fmt.Println(longURL, err)
		case PUT:
			// create a new endpoint
		}
	default:
		// return a 501 error
	}
}

func main() {
	fmt.Println("running")

	http.HandleFunc("/", router)
	http.ListenAndServe(":6000", nil)
}
