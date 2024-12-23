package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	PORT      = "8000"
	BASE_PATH = "https://solace-outputs.s3.ap-south-1.amazonaws.com/__outputs/"
)

func main() {
	r := gin.Default()
	r.Any("/*proxyPath", handleRequest)

	log.Printf("Reverse Proxy Running on port %s..", PORT)
	log.Fatal(r.Run(":" + PORT))
}

func handleRequest(c *gin.Context) {
	hostname := c.Request.Host
	subdomain := strings.Split(hostname, ".")[0]

	resolvesTo := BASE_PATH + "/" + subdomain
	proxyURL, _ := url.Parse(resolvesTo)

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	c.Request.URL.Path = singleJoiningSlash(proxyURL.Path, c.Request.URL.Path)
	if c.Request.URL.Path == "/" {
		c.Request.URL.Path = "/index.html"
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func singleJoiningSlash(a, b string) string {
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
