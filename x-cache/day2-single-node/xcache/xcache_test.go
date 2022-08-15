package xcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestCache_Get(t *testing.T) {
	var f = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("callback failed")
	}
}

// mocker db
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	// 记录查db次数
	loadBbCount := make(map[string]int)
	xCahe := NewGroup("score", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[slowDB]未命中缓存,从db回写数据")
		if v, ok := db[key]; ok {
			if _, ok := loadBbCount[key]; !ok {
				loadBbCount[key] = 0
			}
			loadBbCount[key] += 1

			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		// 首次走db回写
		if view, err := xCahe.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get %s", k)
		}

		// 后续直接走缓存
		if _, err := xCahe.Get(k); err != nil || loadBbCount[k] > 1 {
			t.Fatalf("failed to get %s from key", k)
		}
	}

	if _, err := xCahe.Get("unknow"); err == nil {
		t.Fatalf("key 不存在")
	}
}
