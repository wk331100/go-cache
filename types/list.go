package types

import (
	"sync"
	"time"
)

// NewLists 创建List类型实例
func NewLists() *Lists {
	return &Lists{
		items: make(map[string]*List),
	}
}

// Lists 类型数据结构
type Lists struct {
	mu    sync.Mutex
	items map[string]*List
}

// Exist 判断一个key是否存在
func (ls *Lists) Exist(k string) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	return ls.exist(k)
}

// exist 判断一个key是否存在
func (ls *Lists) exist(k string) bool {
	l, exist := ls.items[k]
	if !exist {
		return false
	}
	if l.isExpired() {
		ls.Del(k)
		return false
	}
	return true
}

// LPush 从队列k的头部，添加一个元素v
func (ls *Lists) LPush(k string, v any) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	l, exist := ls.items[k]
	if !exist {
		l = newList()
	} else if l.isExpired() {
		ls.Del(k)
		return false
	}
	l.LPush(v)
	ls.items[k] = l
	return exist
}

// LPop 从队列k的头部，弹出一个元素
func (ls *Lists) LPop(k string) (any, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	} else if l.isExpired() {
		ls.Del(k)
		return nil, ErrKeyNotExist
	}
	return l.LPop()
}

// RPush 从队列k的尾部，添加一个元素
func (ls *Lists) RPush(k string, v any) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		l = newList()
	} else if l.isExpired() {
		ls.Del(k)
		return false
	}
	l.RPush(v)
	ls.items[k] = l
	return exist
}

// RPop 从队列k的尾部，弹出一个元素
func (ls *Lists) RPop(k string) (any, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	} else if l.isExpired() {
		ls.Del(k)
		return nil, ErrKeyNotExist
	}
	return l.RPop()
}

// LLen 获取队列k的长度
func (ls *Lists) LLen(k string) int {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		return 0
	} else if l.isExpired() {
		ls.Del(k)
		return 0
	}
	return l.LLen()
}

// LRange 获取队列元素列表
func (ls *Lists) LRange(k string, start, stop int) ([]any, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if start > stop {
		return nil, ErrStartStop
	}
	l, exist := ls.items[k]
	if !exist {
		return nil, ErrKeyNotExist
	} else if l.isExpired() {
		ls.Del(k)
		return nil, ErrKeyNotExist
	}
	return l.LRange(start, stop)
}

// Del 删除一个key
func (ls *Lists) Del(k string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	delete(ls.items, k)
}

// Expiration 设置超时时间
func (ls *Lists) Expiration(k string, d time.Duration) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if !ls.exist(k) {
		return ErrKeyNotExist
	}
	ls.items[k].expiration = time.Now().Add(d).UnixNano()
	return nil
}

// ClearExpiration 清理过期的key
func (ls *Lists) ClearExpiration() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for key, item := range ls.items {
		if item.isExpired() {
			delete(ls.items, key)
		}
	}
}

// RandomClearExpiration 随机清理过期的key
func (ls *Lists) RandomClearExpiration() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	var counter int
	for key, item := range ls.items {
		if counter > DefaultCleanItems {
			return
		}
		if item.isExpired() {
			delete(ls.items, key)
		}
		counter++
	}
}

// Flush 清空缓存
func (ls *Lists) Flush() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.items = make(map[string]*List)
}

// newList 创建一个列表的实例
func newList() *List {
	return &List{
		items:      make([]any, 0, 0),
		expiration: DefaultExpiration,
	}
}

// List 列表集合
type List struct {
	items      []any
	expiration int64
}

// LPush 从队列的头部，添加一个元素v
func (l *List) LPush(v any) {
	l.items = append([]any{v}, l.items...)
}

// LPop 从队列的头部，弹出一个元素
func (l *List) LPop() (any, error) {
	var v any
	v, l.items = l.items[0], l.items[1:]
	return v, nil
}

// RPush 从队列的尾部，添加一个元素
func (l *List) RPush(v any) {
	l.items = append(l.items, v)
}

// RPop 从队列的尾部，弹出一个元素
func (l *List) RPop() (any, error) {
	var v any
	v, l.items = l.items[len(l.items)-1], l.items[:len(l.items)-1]
	return v, nil
}

// LLen 获取队列的长度
func (l *List) LLen() int {
	return len(l.items)
}

// LRange 获取队列元素列表
func (l *List) LRange(start, stop int) ([]any, error) {
	if start < 0 {
		start = len(l.items) + start
	}
	if stop < 0 {
		stop = len(l.items) + stop
	}
	if start > stop || start >= len(l.items) {
		return nil, ErrStartStop
	}
	if stop >= len(l.items) {
		stop = len(l.items) - 1
	}
	return l.items[start : stop+1], nil
}

// isExpired 判断一个元素是否过期
func (l *List) isExpired() bool {
	if time.Now().UnixNano() > l.expiration {
		return true
	}
	return false
}
