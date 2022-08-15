package lru

import (
	"container/list"
)

/*
关注点：
1. map + 双向链表；
2. 双向链表冗余 key,方便移除旧缓存时，直接删除缓存；
3. 操作的同时都要更新缓存的大小，用来作为淘汰缓存的依据
 */

type Cache struct {
	maxBytes  int64
	nBytes    int64
	cache     map[string]*list.Element
	ll        *list.List
	onEvicted func(key string, value Value) // 缓存被移除的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		cache:     make(map[string]*list.Element),
		ll:        list.New(),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	// 1.更新缓存
	if ele, ok := c.cache[key]; ok {
		// 移到队列尾
		c.ll.MoveToBack(ele)
		// 获取原值
		kv := ele.Value.(entry)
		// 更新缓存大小
		c.nBytes += value.Len() - kv.value.Len()
		// 更新缓存值
		kv.value = value

	} else { // 2.新增缓存

		// 放到队尾
		ele := c.ll.PushBack(entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.nBytes += value.Len() + int64(len(key)) // 注意：key所占的存储大小也要记录
	}

	// 加数据更新数据导致缓存超限制，引起淘汰缓存
	for c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 移动到队尾
		c.ll.MoveToBack(ele)
		kv := ele.Value.(entry)
		return kv.value, true
	}
	return

}

// 淘汰缓存
func (c *Cache) RemoveOldest() {
	// 获取队首元素
	ele := c.ll.Front()

	if ele == nil {
		return
	}

	// 移除队首元素
	c.ll.Remove(ele)
	kv := ele.Value.(entry)

	// 删除缓存并更新大小
	delete(c.cache, kv.key)
	c.nBytes -= int64(len(kv.key)) + kv.value.Len()

	// 回调函数
	if c.onEvicted != nil {
		c.onEvicted(kv.key, kv.value)
	}
}

// 返回缓存的数量
func (c *Cache) Len() int {
	return c.ll.Len()
}
