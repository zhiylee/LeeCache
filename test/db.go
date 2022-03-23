package main

import (
	"fmt"
	"log"
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
