package main

import (
	"encoding/json"
	"fmt"
	"log"
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
func (a apiHandler) returnJson(resp http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshaling data, %s", err.Error())
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonData)
}

// Given a short URL find the full length URL and returns it
func (a apiHandler) forward(req http.Request, shortURL string) {
	longURL, err := shrink.RetrieveUrl(a.conn, shortURL)
	if err != nil {
		http.Error(a.resp, err.Error(), http.StatusBadRequest)
		log.Printf("%d: URL for short, %s, not found\n", http.StatusFound, shortURL)
		return
	}

	err = shrink.IncrUrlHitCount(a.conn, a.shortURL)
	if err != nil {
		log.Printf(err.Error())
	}

	http.Redirect(a.resp, &req, longURL, 302)
}

// Given a ServerRequest, tries to create a short url before returning
func (a apiHandler) createShortUrl(resp http.ResponseWriter, data ServerRequest) {
	shortURL, err := shrink.CreateURL(a.conn, data.LongURL, data.RequestedURL)
	if err != nil {
		if err == shrink.UrlInUse {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			log.Printf("Short url, %s, already in use\n", shortURL)
			return
		}
		resp.WriteHeader(http.StatusInternalServerError)
		log.Printf("%d: Unclassified error, %s\n", http.StatusInternalServerError, err.Error())
		return
	}
	response := ServerResponse{shortURL}
	a.returnJson(resp, response)
}

func (a apiHandler) getUrlStats(resp http.ResponseWriter, req *http.Request) {
	stats, err := shrink.RetrieveUrlStats(a.conn, a.shortURL)
	if err == shrink.UrlNotFound {
		http.Error(resp, fmt.Sprintf("URL: /%s does not exist\n", a.shortURL), http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("Error retrieving stats, %s\n", err.Error())
		http.Error(resp, fmt.Sprintf("Error retrieving stats for /%s\n", a.shortURL), http.StatusBadRequest)
		return
	}

	statsJson, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error encoding stats json, %s\n", err.Error())
		http.Error(resp, "Error encoding json", http.StatusInternalServerError)
		return
	}

	resp.Write(statsJson)
}
