package lru

import (
	"container/list"
)

type Cache struct {
	MaxEntries int
	OnEvicted func(key Key,value interface{})
	ll *list.List
	cache map[interface{}]*list.Element
}

type Key interface {}

type entry struct {
	key Key
	value interface{}
}
//创建一个新的cache
func New(maxEntrues int)*Cache{
	return &Cache{
		MaxEntries:maxEntrues,
		ll:list.New(),
		cache:make(map[interface{}]*list.Element),
	}
}
//新增一个值到cache
func (c *Cache)Add(key Key,value interface{}){
	if c.cache == nil{
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee,ok := c.cache[key];ok{
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	ele:= c.ll.PushFront(&entry{key,value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries{
		c.RemoveOldest()
	}
}
// Get looks up a key's value from the cache.
func (c *Cache)Get(key Key)(value interface{},ok bool){
	if c.cache == nil{
		return
	}
	if ele,hit := c.cache[key];hit{
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value,true
	}
	return
}
// Remove removes the provided key from the cache.
func (c *Cache)Remove(key Key){
	if c.cache == nil{
		return
	}
	if ele,hit := c.cache[key];hit{
		c.RemoveElement(ele)
	}
}

func (c *Cache)RemoveOldest(){
	if c.cache == nil{
		return
	}
	ele:= c.ll.Back()
	if ele != nil{
		c.RemoveElement(ele)
	}
}
// RemoveOldest removes the oldest item from the cache.
func (c *Cache)RemoveElement(e *list.Element){
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache,kv.key)
	if c.OnEvicted !=nil{
		c.OnEvicted(kv.key,kv.value)
	}
}
// Len returns the number of items in the cache
func (c *Cache)Len()int{
	if c.cache == nil{
		return 0
	}
	return c.ll.Len()
}




