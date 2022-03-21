package leecache

import (
	"fmt"
	"log"
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

func TestGroup(t *testing.T) {
	g := NewGroup("test", GetterFunc(testGetterFunc), 2<<10) //2 << 10 = 2048
	for k, v := range testdb.data {
		//第一次缓存为空，通过回调函数获取数据
		if view, err := g.Get(k); err != nil || view.toString() != v {
			t.Fatalf("failed to get %s=%s, cache value is %s\n", k, v, view.toString())
		}

		//第二次测试是否直接从缓存获取数据，判断缓存是否命中
		//searchCount是对 testDb 查询次数的计数，值 >1 说明缓存未命中
		if view, err := g.Get(k); err != nil || testdb.searchCount[k] > 1 || view.toString() != v {
			t.Fatalf("cache %s=%s miss, cache value is %s", k, v, view.toString())
		}
	}

	//测试get一个不存在的key
	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("get a unkonwn key, but %s got", view.toString())
	}
}
