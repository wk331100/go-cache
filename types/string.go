package types

import (
	"sync"
	"time"
)

// NewStrings 创建字符串类型实例
func NewStrings() *Strings {
	return &Strings{
		items: make(map[string]*Item),
	}
}

// Strings string类型数据结构
type Strings struct {
	mu    sync.Mutex
	items map[string]*Item
}

// Set 设置一个字符串类型
func (s *Strings) Set(k string, v any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.Set(v)
	s.items[k] = i
}

// SetEx 缓存k的值为v,并且设置超时时间d
func (s *Strings) SetEx(k string, v any, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.SetEx(v, d)
	s.items[k] = i
}

// Get 获取一个string类型值
func (s *Strings) Get(k string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	}
	return i.Get(), nil
}

// Incr 对k计数+1
func (s *Strings) Incr(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.Incr()
	s.items[k] = i
}

// Decr 对k计数-1
func (s *Strings) Decr(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.Decr()
	s.items[k] = i
}

// IncrBy 对k计数+v
func (s *Strings) IncrBy(k string, v int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.IncrBy(v)
	s.items[k] = i
}

// DecrBy 对k计数-v
func (s *Strings) DecrBy(k string, v int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.DecrBy(v)
	s.items[k] = i
}

// Del 删除一个key
func (s *Strings) Del(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, k)
}

// newItem 创建一个字符串存储单元的实例
func newItem() *Item {
	return &Item{
		expiration: DefaultExpiration,
	}
}

// Item 每一条string类型数据
// object 数据的具体内容
// expiration 过期时间
type Item struct {
	object     any
	expiration int64
}

func (i *Item) Set(v any) {
	i.object = v
}

func (i *Item) SetEx(v any, d time.Duration) {
	i.object = v
	i.expiration = time.Now().Add(d).UnixNano()
}

func (i *Item) Get() any {
	return i.object
}

func (i *Item) Incr() {
	if i.object == nil {
		i.object = int64(0)
	}
	num := i.object.(int64)
	num++
	i.object = num
}

func (i *Item) IncrBy(v int64) {
	if i.object == nil {
		i.object = int64(0)
	}
	num := i.object.(int64)
	num += v
	i.object = num
}

func (i *Item) Decr() {
	if i.object == nil {
		i.object = int64(0)
	}
	num := i.object.(int64)
	num--
	i.object = num
}

func (i *Item) DecrBy(v int64) {
	if i.object == nil {
		i.object = int64(0)
	}
	num := i.object.(int64)
	num -= v
	i.object = num
}
