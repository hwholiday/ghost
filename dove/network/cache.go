package network

import (
	"encoding/json"
	"sync"
)

type Cache struct {
	mp  sync.Map
	tmp string
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Save(k string, v any) {
	c.mp.Store(k, v)
}

func (c *Cache) Get(k string) *Cache {
	c.tmp = k
	return c
}
func (c *Cache) load() (any, bool) {
	defer c.clearTmp()
	return c.mp.Load(c.tmp)
}

func (c *Cache) Int() int {
	v, ok := c.load()
	if !ok {
		return 0
	}
	vStr, ok := v.(int)
	if !ok {
		return 0
	}
	return vStr
}
func (c *Cache) clearTmp() {
	c.tmp = ""
}

func (c *Cache) Result() (any, bool) {
	return c.load()
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

func (c *Cache) String() string {
	if c.tmp == "" {
		str, _ := json.Marshal(c.Map())
		return string(str)
	}
	v, ok := c.load()
	if !ok {
		return ""
	}
	vStr, ok := v.(string)
	if !ok {
		return ""
	}
	return vStr
}
