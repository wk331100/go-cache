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

// Exist 判断k是否存在
func (s *Strings) Exist(k string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.exist(k)
}

// exist 判断k是否存在
func (s *Strings) exist(k string) bool {
	i, exist := s.items[k]
	if !exist {
		return false
	}
	if i.isExpired() {
		s.del(k)
		return false
	}
	return true
}

// Set 设置一个字符串类型
func (s *Strings) Set(k string, v any) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.Set(v)
	s.items[k] = i
	return exist
}

// SetEx 缓存k的值为v,并且设置超时时间d
func (s *Strings) SetEx(k string, v any, d time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		i = newItem()
	}
	i.SetEx(v, d)
	s.items[k] = i
	return exist
}

// Get 获取一个string类型值
func (s *Strings) Get(k string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, exist := s.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	} else if i.isExpired() {
		s.del(k)
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
	} else if i.isExpired() {
		s.del(k)
		return
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
	} else if i.isExpired() {
		s.del(k)
		return
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
	} else if i.isExpired() {
		s.del(k)
		return
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
	} else if i.isExpired() {
		s.del(k)
		return
	}
	i.DecrBy(v)
	s.items[k] = i
}

// Del 删除一个key
func (s *Strings) Del(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.del(k)
}

func (s *Strings) del(k string) {
	delete(s.items, k)
}

// Expiration 设置超时时间
func (s *Strings) Expiration(k string, d time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.exist(k) {
		return ErrKeyNotExist
	}
	s.items[k].expiration = time.Now().Add(d).UnixNano()
	return nil
}

// ClearExpiration 清理过期的key
func (s *Strings) ClearExpiration() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, item := range s.items {
		if item.isExpired() {
			delete(s.items, key)
		}
	}
}

// RandomClearExpiration 随机清理100条过期的key
func (s *Strings) RandomClearExpiration() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var counter int
	for key, item := range s.items {
		if counter > DefaultCleanItems {
			return
		}
		if item.isExpired() {
			delete(s.items, key)
		}
		counter++
	}
}

// Flush 清空缓存
func (s *Strings) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[string]*Item)
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

// isExpired 判断一个元素是否过期
func (i *Item) isExpired() bool {
	if i.expiration != DefaultExpiration && time.Now().UnixNano() > i.expiration {
		return true
	}
	return false
}
