package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nowBytes int64
	ll       *list.List // 双向链表
	cache    map[string]*list.Element

	// Todo:callback on add && delete
}

// 缓存的对象 双向链表的数据类型
type entry struct {
	key   string
	value Value
}

func (ent entry) len() int {
	return len(ent.key) + ent.value.len()
}

// 缓存的数据
type Value interface {
	len() int
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}

}
