package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	PORT      = "8000"                                                               // Port number where the server will run
	BASE_PATH = "https://vercel-clone-outputs.s3.ap-south-1.amazonaws.com/__outputs" // Base URL path for resolving subdomains
)

func main() {
	r := gin.Default()                  // Create a Gin router with default middleware: logger and recovery
	r.Any("/*proxyPath", handleRequest) // Route all requests to handleRequest function

	log.Printf("Reverse Proxy Running on port %s..", PORT) // Log message
	log.Fatal(r.Run(":" + PORT))                           // Run the server on the specified port
}

func handleRequest(c *gin.Context) {
	hostname := c.Request.Host                   // Extract the hostname from the incoming request
	subdomain := strings.Split(hostname, ".")[0] // Get the subdomain (e.g., if hostname is "sub.example.com", subdomain will be "sub")

	// Custom Domain - DB Query
	// Example: Assume the hostname is "myapp.example.com"
	// Subdomain extracted: "myapp"

	resolvesTo := BASE_PATH + "/" + subdomain // Construct the URL to which the request should be proxied
	// Example: resolvesTo = "https://vercel-clone-outputs.s3.ap-south-1.amazonaws.com/__outputs/myapp"

	proxyURL, _ := url.Parse(resolvesTo) // Parse the target URL

	proxy := httputil.NewSingleHostReverseProxy(proxyURL) // Create a reverse proxy pointing to the target URL

	// Adjust the URL path for the proxy request
	c.Request.URL.Path = singleJoiningSlash(proxyURL.Path, c.Request.URL.Path)
	// Example: If incoming request URL is "/", it will be changed to "/index.html"

	if c.Request.URL.Path == "/" {
		c.Request.URL.Path = "/index.html"
	}

	proxy.ServeHTTP(c.Writer, c.Request) // Forward the request to the target URL
}

func singleJoiningSlash(a, b string) string {
	// Helper function to join URL paths correctly
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
