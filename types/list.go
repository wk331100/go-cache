package types

import "sync"

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

// LPush 从队列k的头部，添加一个元素v
func (ls *Lists) LPush(k string, v any) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	l, exist := ls.items[k]
	if !exist {
		l = newList()
	}
	l.LPush(v)
	ls.items[k] = l
}

// LPop 从队列k的头部，弹出一个元素
func (ls *Lists) LPop(k string) (any, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		return nil, ErrEmptyList
	}
	return l.LPop()
}

// RPush 从队列k的尾部，添加一个元素
func (ls *Lists) RPush(k string, v any) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		l = newList()
	}
	l.RPush(v)
	ls.items[k] = l
}

// RPop 从队列k的尾部，弹出一个元素
func (ls *Lists) RPop(k string) (any, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	l, exist := ls.items[k]
	if !exist {
		return nil, ErrEmptyList
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
		return nil, ErrEmptyList
	}
	return l.LRange(start, stop)
}

// Del 删除一个key
func (ls *Lists) Del(k string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	delete(ls.items, k)
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
