package main

import (
	"encoding/json"
	"net/http"
	"github.com/natebrennand/shrtnr/shrink"
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
func (a apiHandler) ReturnJson(resp http.ResponseWriter, data ServerResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonData)
}

// Given a short URL find the full length URL and returns it
func (a apiHandler) getLongUrl(req http.Request, shortURL string) {
	longURL, err := shrink.RetrieveURL(a.conn, shortURL)
	if err != nil {
		if err == shrink.UrlNotFound {
			http.Error(a.resp, err.Error(), http.StatusFound)
			return
		}
		http.Error(a.resp, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(a.resp, &req, longURL, 302)
}

// Given a ServerRequest, tries to create a short url before returning
func (a apiHandler) createShortUrl(resp http.ResponseWriter, data ServerRequest) {
	shortURL, err := shrink.CreateURL(a.conn, data.LongURL, data.RequestedURL)
	if err != nil {
		if err == shrink.UrlInUse {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := ServerResponse{shortURL}
	a.ReturnJson(resp, response)
}
