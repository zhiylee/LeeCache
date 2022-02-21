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
	c.nowBytes -= int64(ent.len())

}

// Add or Update entry
func (c *Cache) Add(key string, value Value) bool {
	if ele, ok := c.cache[key]; ok { //尝试获取map，存在则更新
		c.ll.MoveToFront(ele)
		ent := ele.Value.(*entry)
		c.nowBytes += int64(value.len() - ent.value.len()) //更新内存
		ent.value = value
	} else { //不存在则新增
		ent := entry{key, value}
		ele := c.ll.PushFront(&ent)
		c.cache[key] = ele
		c.nowBytes += int64(ent.len())
	}

	//对超出的内存进行回收
	//先增添再回收，占用的内存会少量溢出maxBytes
	for c.maxBytes != 0 && c.nowBytes != 0 && c.nowBytes > c.maxBytes {
		c.Eliminate()
	}

	return true
}

// Len of cache entry
func (c *Cache) Len() int {
	return c.ll.Len()
}
