package leecache

import (
	"testing"
)

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
