package cache

import (
	"container/list"
	"sync"

	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/model"
)

type Cache struct {
	mu    sync.RWMutex
	data  map[string]*list.Element
	order *list.List // keep insertion order (FIFO)
	limit int
	log   logger.InterfaceLogger
}

type entry struct {
	key   string
	value *model.Order
}

var _ InterfaceCache = (*Cache)(nil)

func NewCache(log logger.InterfaceLogger) *Cache {
	return &Cache{
		data:  make(map[string]*list.Element),
		order: list.New(),
		limit: 10,
		log:   log,
	}
}

func (c *Cache) Get(key string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.data[key]
	if !ok {
		c.log.Infof("Key not found: %s", key)
		return nil, false
	}

	ent := elem.Value.(*entry)
	c.log.Infof("Get from cache: %s", key)
	return ent.value, true
}

func (c *Cache) Set(key string, value *model.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if already exists, update
	if elem, ok := c.data[key]; ok {
		c.log.Infof("Update in cache: %s", key)
		elem.Value.(*entry).value = value
		return nil
	}

	// check limit
	if c.order.Len() >= c.limit {
		// remove oldest
		oldest := c.order.Front()
		if oldest != nil {
			ent := oldest.Value.(*entry)
			delete(c.data, ent.key)
			c.order.Remove(oldest)
			c.log.Infof("Removed oldest from cache: %s", ent.key)
		}
	}

	ent := &entry{key, value}
	elem := c.order.PushBack(ent)
	c.data[key] = elem
	c.log.Infof("Set to cache: %s", key)
	return nil
}
