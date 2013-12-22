package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shrtnr/shrink"
)

const (
	GET  string = "GET"
	POST string = "POST"
	PUT  string = "PUT"
)

// Used for server responses
type ServerResponse struct {
	URL string
}

// Used for server requests
type ServerRequest struct {
	LongURL      string
	RequestedURL string
}

// Serializes and returns the given ServerResponse struct through the resp
func ReturnJson(resp http.ResponseWriter, data ServerResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonData)
}

// Parses the JSON request body from the request
func GetReqBody(req http.Request) (ServerRequest, error) {
	decoder := json.NewDecoder(req.Body)
	var requestBody ServerRequest
	err := decoder.Decode(&requestBody)
	return requestBody, err
}

// Given a short URL find the full length URL and returns it
func GetFullURL(resp http.ResponseWriter, req http.Request, shortURL string) {
	longURL, err := shrink.RetrieveURL(shortURL)
	if err != nil {
		if err == shrink.UrlNotFound {
			http.Error(resp, err.Error(), http.StatusFound)
			return
		}
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(resp, &req, longURL, 302)
}

// Given a ServerRequest, tries to create a short url before returning
func CreateURL(resp http.ResponseWriter, data ServerRequest) {
	shortURL, err := shrink.CreateURL(data.LongURL, data.RequestedURL)
	if err != nil {
		if err == shrink.UrlInUse {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := ServerResponse{shortURL}
	fmt.Println(response)
	ReturnJson(resp, response)
}

// Routes all http requests to their corresponding functions
func router(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	method := req.Method
	shortURL := url[1:]

	fmt.Println(method, url) // TODO: delete

	switch {
	case method == GET:
		switch url {
		case "/":
			// return homepage
			// TODO: deal w/ static resources too
		default:
			GetFullURL(resp, *req, shortURL)
		}
	case method == POST && url == "/":
		requestBody, err := GetReqBody(*req)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
		}
		CreateURL(resp, requestBody)
	default:
		// return a 501 error
		resp.WriteHeader(http.StatusNotImplemented)
	}
}

func main() {
	fmt.Println("running on 6000")

	http.HandleFunc("/", router)
	http.ListenAndServe(":6000", nil)
}
