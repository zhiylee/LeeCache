package lru

import (
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache_Get(t *testing.T) {
	lru := New(1000)
	lru.Add("k1", String("value1"))
	if v, ok := lru.Get("k1"); !ok || string(v.(String)) != "value1" {
		t.Fatalf("Get cache k1=value1 failed")
	}
	lru.Add("k2", String("value2"))
	if v, ok := lru.Get("k2"); !ok || string(v.(String)) != "value2" {
		t.Fatalf("Get cache k2=value2 failed")
	}
}

func TestCache_Eliminate(t *testing.T) {
	k := []string{"k1", "k2", "k3"}
	v := []string{"value1", "value2", "value3"}

	lru := New(int64(len(k[0] + k[1] + v[0] + v[1])))
	lru.Add(k[0], String(v[0]))
	lru.Add(k[1], String(v[1]))
	lru.Add(k[2], String(v[2]))

	if _, ok := lru.Get(k[0]); ok || lru.Len() != 2 {
		t.Fatalf("Eliminate failed,not deleted the oldest entry k1")
	}

	if value, ok := lru.Get(k[2]); !ok || string(value.(String)) != v[2] {
		t.Fatalf("Add k2=value2 failed")
	}
}
