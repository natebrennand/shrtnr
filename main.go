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

	fmt.Println(method, url)

	switch {
	case method == GET:
		switch url {
		case "/":
			// return homepage
		default:
			// returned full length URL
			longURL, err := shrink.RetrieveURL(shortURL)
			fmt.Println(longURL, err)
		}
	case method == POST && url == "/":
		// creates a random endpoint
	case method == PUT && url != "/":
		// creates a specific endpoint
	default:
		resp.WriteHeader(http.StatusNotImplemented)
		// return a 501 error
	}
}

func main() {
	fmt.Println("running on 6000")

	http.HandleFunc("/", router)
	http.ListenAndServe(":6000", nil)
}
