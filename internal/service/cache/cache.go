package cache

import (
	"github.com/merkulovlad/wbtech-go/internal/logger"
	"github.com/merkulovlad/wbtech-go/internal/model"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]*model.Order
	log  logger.InterfaceLogger
}

var _ InterfaceCache = (*Cache)(nil)

func NewCache(log logger.InterfaceLogger) *Cache {
	return &Cache{
		data: make(map[string]*model.Order),
		log:  log,
		mu:   sync.RWMutex{},
	}
}

func (c *Cache) Get(key string) (*model.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	order, ok := c.data[key]
	if !ok {
		c.log.Infof("Key not found: %s", key)
		return nil, false
	}
	c.log.Infof("Get from cache: %s", key)
	return order, true
}

func (c *Cache) Set(key string, value *model.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	c.log.Infof("Set to cache: %s", key)
	return nil
}
