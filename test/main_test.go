package main

import (
	"LeeCache/peers"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGet(t *testing.T, url string) {
	res, err := http.Get(url)
	if err != nil {
		t.Log("get fail : ", err)
		return
	}
	defer res.Body.Close()

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Log("get value faid")
		return
	}

	if string(value) != "va" {
		t.Fatalf("get value error %v", value)
	}

	wg.Done()
}

var wg sync.WaitGroup

func TestPressure(t *testing.T) {

	n := 10000

	wg.Add(n)

	url := "http://127.0.0.1:8001/api?group=test&key=a"
	for i := 0; i < n; i++ {
		go httpGet(t, url)
		if i%100 == 0 {
			time.Sleep(time.Millisecond * 20)
		}
	}

	wg.Wait()

}

func get(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Println("get fail : ", err)
		return
	}
	defer res.Body.Close()

	value, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("get value faid")
		return
	}

	if string(value) != "va" {
		log.Printf("get value error %v", value)
	}
}

func BenchmarkGetPeer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "http://127.0.0.1:8002/api?group=test&key=a"
		get(url)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//g := peers.NewHTTPGetter("localhost:9981")
		g := peers.NewRPCGetter("localhost:9991")
		v, err := g.Get("test", "a")
		if err != nil {
			b.Fatal("err")
		}
		if string(v) != "va" {
			b.Fatal("value error ", v)
		}

	}
}
