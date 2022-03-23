package main

import (
	leecache "LeeCache"
	"LeeCache/peers"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func startHTTPServe(addr string) {
	serve := leecache.NewHTTPServe(addr)
	log.Printf("[LeeCache] HTTP serve is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr, serve))
}

func main() {
	var id int
	var replicas int
	flag.IntVar(&id, "id", 1, "LeeCache peer id")
	flag.IntVar(&replicas, "replicas", 50, "the number of peer's replicas")
	flag.Parse()

	pool := peers.NewPool(replicas)
	leecache.Localhost = "localhost:" + "900" + strconv.Itoa(id)

	addrMap := map[string]string{
		"Peer1": "localhost:9001",
		"Peer2": "localhost:9002",
		"Peer3": "localhost:9003",
	}

	for name, addr := range addrMap {
		pool.Add(peers.NewPeer(name, addr))
	}

	g := leecache.NewGroup("test", leecache.GetterFunc(testGetterFunc), 2<<10)
	g.Peers = pool

	go startApiServe(":800" + strconv.Itoa(id))
	startHTTPServe(":900" + strconv.Itoa(id))
}
