package leecache

import (
	"LeeCache/peers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"
)

type testDb struct {
	data        map[string]string
	searchCount map[string]int
}

//获取数据并对获取的次数进行计数
func (d testDb) Get(key string) (string, error) {
	if _, ok := d.searchCount[key]; !ok {
		d.searchCount[key] = 1
	} else {
		d.searchCount[key]++
	}
	log.Printf("[testDb] search key %s, times %d", key, d.searchCount[key])
	if v, ok := d.data[key]; !ok {
		return "", fmt.Errorf("key %s not exist", key)
	} else {
		return v, nil
	}
}

var testData = map[string]string{
	"a": "va",
	"b": "vb",
	"c": "vc",
}

var testdb = testDb{
	data:        testData,
	searchCount: make(map[string]int),
}

func testGetterFunc(key string) ([]byte, error) {
	v, err := testdb.Get(key)

	return []byte(v), err
}

func startHTTPServe(addr string) {
	serve := NewHTTPServe(addr)
	log.Printf("[LeeCache] HTTP serve is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, serve))
}

func Test(t *testing.T) {
	var port int
	flag.IntVar(&port, "port", 8001, "LeeCache HTTP server port")
	flag.Parse()

	pool := peers.NewPool(50)

	addrMap := map[string]string{
		"Peer1": "http://localhost:8001",
		"Peer2": "http://localhost:8002",
		"Peer3": "http://localhost:8003",
	}

	for name, addr := range addrMap {
		pool.Add(peers.NewPeer(name, addr))
	}

	g := NewGroup("test", GetterFunc(testGetterFunc), 2<<10)
	g.Peers = pool
	startHTTPServe(":" + strconv.Itoa(port))
}
