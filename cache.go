package go_cache

import (
	"sync"
	"time"
)

const (
	DefaultExpiration = -1
)

// item 每一条string类型数据
// object 数据的具体内容
// expiration 过期时间
type item struct {
	object     any
	expiration int64
}

// hashItem 每一个Hash类型数据
// object hash内的field内容
// expiration 过期时间
type hashItem struct {
	object     map[string]any
	expiration int64
}

// NewCache 创建新的缓存对象
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]item),
		lists: make(map[string][]any),
		hash:  make(map[string]hashItem),
	}
}

// Cache 缓存结构
// items 为string类型的映射
type Cache struct {
	itemLocker sync.RWMutex
	listLocker sync.Mutex
	hashLocker sync.RWMutex
	items      map[string]item
	lists      map[string][]any
	hash       map[string]hashItem
	sets       map[string][]any
}

// ======== 字符串 =======

// Set 缓存k的值为v
func (c *Cache) Set(k string, v any) error {
	c.itemLocker.Lock()
	defer c.itemLocker.Unlock()
	c.items[k] = item{
		object:     v,
		expiration: DefaultExpiration,
	}
	return nil
}

// SetEx 缓存k的值为v,并且设置超时时间d
func (c *Cache) SetEx(k string, v any, d time.Duration) error {
	c.itemLocker.Lock()
	defer c.itemLocker.Unlock()
	c.items[k] = item{
		object:     v,
		expiration: time.Now().Add(d).UnixNano(),
	}
	return nil
}

// Get 获取一个string类型值
func (c *Cache) Get(k string) (any, error) {
	c.itemLocker.RLock()
	c.itemLocker.RUnlock()
	v, exist := c.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	}
	return v.object, nil
}

// Incr 对k计数+1
func (c *Cache) Incr(k string) {
	c.itemLocker.Lock()
	c.itemLocker.Unlock()
	i, exist := c.items[k]
	if !exist {
		i = item{
			object:     int64(0),
			expiration: DefaultExpiration,
		}
	}
	num := i.object.(int64)
	i.object = num + 1
	c.items[k] = i
}

// Decr 对k计数-1
func (c *Cache) Decr(k string) {
	c.itemLocker.Lock()
	c.itemLocker.Unlock()
	i, exist := c.items[k]
	if !exist {
		i = item{
			object:     int64(0),
			expiration: DefaultExpiration,
		}
	}
	num := i.object.(int64)
	i.object = num - 1
	c.items[k] = i
}

// IncrBy 对k计数+v
func (c *Cache) IncrBy(k string, v int64) {
	c.itemLocker.Lock()
	c.itemLocker.Unlock()
	i, exist := c.items[k]
	if !exist {
		i = item{
			object:     int64(0),
			expiration: DefaultExpiration,
		}
	}
	num := i.object.(int64)
	i.object = num + v
	c.items[k] = i
}

// DecrBy 对k计数-v
func (c *Cache) DecrBy(k string, v int64) {
	c.itemLocker.Lock()
	c.itemLocker.Unlock()
	i, exist := c.items[k]
	if !exist {
		i = item{
			object:     int64(0),
			expiration: DefaultExpiration,
		}
	}
	num := i.object.(int64)
	i.object = num - v
	c.items[k] = i
}

// ======== 列表 =======

// LPush 从队列k的头部，添加一个元素v
func (c *Cache) LPush(k string, v any) {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	l, exist := c.lists[k]
	if !exist {
		l = make([]any, 0)
	}
	c.lists[k] = append([]any{v}, l...)
}

// LPop 从队列k的头部，弹出一个元素
func (c *Cache) LPop(k string) (any, error) {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	l, exist := c.lists[k]
	if !exist {
		return nil, ErrEmptyList
	}
	v, l := l[0], l[1:]
	c.lists[k] = l
	return v, nil
}

// RPush 从队列k的尾部，添加一个元素
func (c *Cache) RPush(k string, v any) {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	l, exist := c.lists[k]
	if !exist {
		l = make([]any, 0)
	}
	c.lists[k] = append(l, v)
}

// RPop 从队列k的尾部，弹出一个元素
func (c *Cache) RPop(k string) (any, error) {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	l, exist := c.lists[k]
	if !exist {
		return nil, ErrEmptyList
	}
	v, l := l[len(l)-1], l[0:len(l)-1]
	c.lists[k] = l
	return v, nil
}

// LLen 获取队列k的长度
func (c *Cache) LLen(k string) int {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	l, exist := c.lists[k]
	if !exist {
		return 0
	}
	return len(l)
}

// LRange 获取队列元素列表
func (c *Cache) LRange(k string, start, len int) ([]any, error) {
	c.listLocker.Lock()
	c.listLocker.Unlock()
	if start < 0 || len <= 0 {
		return nil, ErrStartLen
	}

	l, exist := c.lists[k]
	if !exist {
		return nil, ErrEmptyList
	}
	var vList []any
	for index, v := range l {
		if index >= start && index < start+len {
			vList = append(vList, v)
		}
	}
	return vList, nil
}

// ======== 散列Hash =======

// HSet 缓存数据到Hash中
// k 为Hash中的key
// field 为hash中项
func (c *Cache) HSet(k, field string, v any) {
	c.hashLocker.Lock()
	defer c.hashLocker.Unlock()
	h, exist := c.hash[k]
	if !exist {
		h = hashItem{
			object:     make(map[string]any),
			expiration: -1,
		}
		c.hash[k] = h
	}
	c.hash[k].object[field] = v
}

// HGet 从Hash中获取存储的元素
func (c *Cache) HGet(k, field string) (any, error) {
	c.hashLocker.Lock()
	defer c.hashLocker.Unlock()
	h, exist := c.hash[k]
	if !exist {
		return nil, ErrHashKey
	}
	v, exist := h.object[field]
	if !exist {
		return nil, ErrHashField
	}
	return v, nil
}

func (c *Cache) HDel(k, field string) {
	c.hashLocker.Lock()
	defer c.hashLocker.Unlock()
	_, exist := c.hash[k]
	if exist {
		delete(c.hash[k].object, field)
	}
}

func (c *Cache) HKeys(k string) ([]string, error) {
	c.hashLocker.Lock()
	defer c.hashLocker.Unlock()
	_, exist := c.hash[k]
	if !exist {
		return nil, ErrHashKey
	}
	var fields []string
	for field, _ := range c.hash[k].object {
		fields = append(fields, field)
	}
	return fields, nil
}

func (c *Cache) HVals(k string) ([]any, error) {
	c.hashLocker.Lock()
	defer c.hashLocker.Unlock()
	_, exist := c.hash[k]
	if !exist {
		return nil, ErrHashKey
	}
	var vals []any
	for _, v := range c.hash[k].object {
		vals = append(vals, v)
	}
	return vals, nil
}

// ======== 集合 =======

func (c *Cache) SAdd(k string, m any) error {
	return nil
}

func (c *Cache) SRem(k, m string) error {
	return nil
}

func (c *Cache) SMembers(k string) ([]any, error) {
	return nil, nil
}

func (c *Cache) SIsMember(k string, m any) bool {
	return false
}

func (c *Cache) SCard(k string, m any) int64 {
	return 0
}

// ======== 全局 =======

// Del 删除一个key
func (c *Cache) Del(k string) {

}

// Expiration 设置超时时间
func (c *Cache) Expiration(k string, d time.Duration) error {
	return nil
}
