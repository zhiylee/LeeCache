package lru

// Get the value and reset lifespan
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ent := ele.Value.(*entry)
		return ent.value, true
	}

	return
}

// Eliminate the oldest entry
func (c *Cache) Eliminate() {
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	ent := ele.Value.(*entry)
	c.ll.Remove(ele)
	delete(c.cache, ent.key)
	c.nowBytes -= int64(ent.Len())

}

// Add or Update entry
func (c *Cache) Add(key string, value Value) bool {
	var requireBytes int64 = 0
	if ele, ok := c.cache[key]; ok { //尝试获取map，存在则更新
		c.ll.MoveToFront(ele)
		ent := ele.Value.(*entry)
		requireBytes = int64(value.Len() - ent.value.Len())
		c.recycleSpace(requireBytes) //腾出足够的空间
		ent.value = value
	} else { //不存在则新增
		ent := entry{key, value}
		requireBytes = int64(ent.Len())
		c.recycleSpace(requireBytes)
		ele := c.ll.PushFront(&ent)
		c.cache[key] = ele
	}

	c.nowBytes += requireBytes //更新内存占用大小

	return true
}

//内存回收 腾出足够的空间
func (c *Cache) recycleSpace(requireBytes int64) {
	if requireBytes <= 0 || c.maxBytes <= 0 {
		return
	}
	for c.maxBytes-c.nowBytes < requireBytes {
		c.Eliminate()
	}
}

// Len of cache entry
func (c *Cache) Len() int {
	return c.ll.Len()
}
