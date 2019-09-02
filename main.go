package main

import (
	"fmt"
	"os"
	"log"
	"strings"
	"net/http"
	"net/http/httputil"
	"github.com/bobziuchkovski/digest"
)

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return defaultValue
    }
    return value
}

func splitHostAndPath(path string) (string, string) {
	p := strings.Index(path[1:], "/") + 1
	return path[1:p], path[p:]
}

func NewReverseProxy(username string, password string) *httputil.ReverseProxy {
	transport := digest.NewTransport(username, password)
	director := func(req *http.Request) {
		host, path := splitHostAndPath(req.URL.Path)
		req.URL.Scheme = "https"
		req.URL.Host = host
		req.URL.Path = path
	}
	return &httputil.ReverseProxy{
		Director: director,
		Transport: transport,
	}
}

func main() {
	proxy := NewReverseProxy(getEnv("DIGEST_USERNAME", ""), getEnv("DIGEST_PASSWORD", ""))
	host := fmt.Sprintf("%s:%s", getEnv("HOSTNAME", ""), getEnv("PORT", "80"))
	fmt.Printf("Listen %s\n", host)
	err := http.ListenAndServe(host, proxy)
	if err != nil {
		log.Fatal(err)
	}
}
