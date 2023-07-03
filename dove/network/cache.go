package network

import (
	"encoding/json"
	"sync"
)

type Cache struct {
	mp sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Save(k string, v any) {
	c.mp.Store(k, v)
}

func (c *Cache) GetString(k string) string {
	v, ok := c.Result(k)
	if !ok {
		return ""
	}
	vStr, ok := v.(string)
	if !ok {
		return ""
	}
	return vStr
}

func (c *Cache) GetInt(k string) int {
	v, ok := c.Result(k)
	if !ok {
		return 0
	}
	vInt, ok := v.(int)
	if !ok {
		return 0
	}
	return vInt
}

func (c *Cache) Result(k string) (any, bool) {
	return c.mp.Load(k)
}

func (c *Cache) String() string {
	str, _ := json.Marshal(c.Map())
	return string(str)
}

func (c *Cache) Map() map[string]any {
	var result = make(map[string]any)
	c.mp.Range(func(key, value any) bool {
		if keyStr, ok := key.(string); ok {
			result[keyStr] = value
		}
		return true
	})
	return result
}
