package leecache

import (
	"log"
	"net/http"
	"testing"
)

func TestHTTPServe(t *testing.T) {
	NewGroup("test", GetterFunc(testGetterFunc), 2<<10)

	addr := ":8888"
	serve := NewHTTPServe(addr)
	log.Printf("[LeeCache] HTTP serve is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, serve))
}
