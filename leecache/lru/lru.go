package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nowBytes int64
	ll       *list.List // 双向链表
	cache    map[string]*list.Element

	// Todo:callback on add && delete
}

// 缓存的对象
type entry struct {
	key   string
	value Value
}

func (ent entry) Len() int {
	return len(ent.key) + ent.value.Len()
}

// Value 缓存的数据
type Value interface {
	Len() int
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}

}
