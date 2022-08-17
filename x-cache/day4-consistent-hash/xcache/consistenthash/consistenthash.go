package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/* 新知识点:
1. 函数类型依赖注入，便于替换函数
2. sort.Search() 找到最小满足条件的索引
 */


type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int // Sorted hash 环
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if fn == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(key + strconv.Itoa(i))))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string{
	if m.hashMap == nil{
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] > hash
	})

	// idx == len(m.keys) 取 m.keys[0]
 	return m.hashMap[m.keys[idx%len(m.keys)]]
}
