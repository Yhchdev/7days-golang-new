package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int64 {
	return int64(len(s))
}

// test add and get
func TestCache_Get(t *testing.T) {
	k1, v1 := "key1", "value1"
	lru := NewCache(int64(len(k1)+len(v1)), nil)
	lru.Add("key1", String("value1"))

	// cache hit
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "value1" {
		t.Fatalf("get key1 failed")
	}

	// cache miss
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("get key2 failed")
	}
}

// test remove oldest and update len
func TestCacheLen(t *testing.T) {
	lru := NewCache(int64(6), nil)
	lru.Add("key", String("1"))
	lru.Add("key", String("111"))

	if lru.nBytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nBytes)
	}
}

// test remove oldest
func TestRemoveOldest(t *testing.T) {
	k1, v1 := "key1", "value1"
	k2, v2 := "key2", "value2"
	k3, v3 := "key3", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	// cache miss
	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("remove oldest failed")
	}

}

// test evicted
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callBack := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := NewCache(int64(10), callBack)
	lru.Add("k1", String("123456"))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))

	except := []string{"k1", "k2"}
	if !reflect.DeepEqual(keys, except) {
		t.Fatalf("call onEvicted failed")
	}

}
