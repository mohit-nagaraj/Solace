package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	PORT       = ":8000"
	TARGET_URL = "https://solace-outputs.s3.ap-south-1.amazonaws.com"
)

func NewReverseProxy() (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(TARGET_URL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Get subdomain from the host
		subdomain := strings.Split(req.Host, ".")[0]

		// Keep the original path after the subdomain
		originalPath := req.URL.Path
		if originalPath == "/" {
			originalPath = "/index.html"
		}

		// Set the path using the subdomain
		req.URL.Path = "/__outputs/" + subdomain + originalPath

		// Set the host header to match the target
		req.Host = targetURL.Host
	}

	return proxy, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	proxy, err := NewReverseProxy()
	if err != nil {
		log.Printf("Error creating reverse proxy: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("Reverse Proxy Running on port %s", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
