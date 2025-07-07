package hw04lrucache

type cacheItem struct {
	key   Key
	value interface{}
}

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, exists := c.items[key]; exists {

		item.Value = cacheItem{key: key, value: value}

		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() >= c.capacity {
		back := c.queue.Back()
		if back != nil {
			ci := back.Value.(cacheItem)
			delete(c.items, ci.key)
			c.queue.Remove(back)
		}
	}

	item := c.queue.PushFront(cacheItem{key: key, value: value})
	c.items[key] = item

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
