package leecache

// ByteView 缓存的value
type ByteView struct {
	b []byte //真实的缓存值
}

//Len 缓存的长度 实现 lru.Value 接口
func (v ByteView) Len() int {
	return len(v.b)
}

//ByteCopy 返回一个真实数据的拷贝 对外部而言 v.b 只读
func (v ByteView) ByteCopy() []byte {
	b := make([]byte, v.Len())
	copy(b, v.b)
	return b
}

//toString 将真实数据作为字符串返回
func (v ByteView) toString() string {
	return string(v.b)
}
