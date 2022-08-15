package xcache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache Cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己
// 将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。

// 函数类型
type GetterFunc func(key string) ([]byte, error)

// 接口函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: Cache{cacheBytes: cacheBytes},
	}

	groups[name] = g

	return g

}

func GetGroup(name string) *Group {
	// 只读操作，用读锁就可以了
	mu.RLocker()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 缓存命中
	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[XCache] hit")
		return v, nil
	}

	// 缓存未命中 回写数据
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	values := ByteView{bytes}
	g.populateCache(key, values)
	return values, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value)
}
