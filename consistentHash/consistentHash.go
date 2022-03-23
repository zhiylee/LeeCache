package consistentHash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(b []byte) uint32

type Map struct {
	hash     Hash
	replicas int            //number of virtual nodes
	keys     []int          //Sorted
	hashMap  map[int]string //virtual node map to real
}

func New(replicas int, fn Hash) *Map {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}

	return &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

//Add : add keys to the hash ring
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

// Get : get the closest item in the hash ring
func (m *Map) Get(key string) string {
	if len(m.keys) < 1 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// binary search
	id := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// back to the head of hash ring
	if id == len(m.keys) {
		id = 0
	}

	return m.hashMap[m.keys[id]]
}

// Delete : delete the item in the hash ring
func (m *Map) Delete(key string) {
	// Todo
}
