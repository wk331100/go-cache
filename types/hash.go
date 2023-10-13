package types

import "sync"

// NewHashes 创建Hashes类型实例
func NewHashes() *Hashes {
	return &Hashes{
		items: make(map[string]*Hash),
	}
}

// Hashes Hashes类型数据结构
type Hashes struct {
	mu    sync.Mutex
	items map[string]*Hash
}

// HSet 缓存数据到Hash中
// k 为Hash中的key
// field 为hash中项
func (hs *Hashes) HSet(k, field string, v any) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	h, exist := hs.items[k]
	if !exist {
		h = newHash()
	}
	h.HSet(field, v)
	hs.items[k] = h
}

// HGet 从Hash中获取存储的元素
func (hs *Hashes) HGet(k, field string) (any, error) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	h, exist := hs.items[k]
	if !exist {
		return nil, ErrHashKey
	}
	return h.HGet(field)
}

// HDel 从Hash中删除元素field
func (hs *Hashes) HDel(k, field string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	h, exist := hs.items[k]
	if exist {
		h.HDel(field)
	}
}

// HKeys 获取Hash中的所有元素field
func (hs *Hashes) HKeys(k string) ([]string, error) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	h, exist := hs.items[k]
	if !exist {
		return nil, ErrHashKey
	}
	return h.HKeys()
}

// HVals 获取Hash中所有元素的内容
func (hs *Hashes) HVals(k string) ([]any, error) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	h, exist := hs.items[k]
	if !exist {
		return nil, ErrHashKey
	}
	return h.HVals()
}

// Del 删除一个key
func (hs *Hashes) Del(k string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	delete(hs.items, k)
}

// newHash 创建一个Hash的实例
func newHash() *Hash {
	return &Hash{
		fields:     make(map[string]any),
		expiration: DefaultExpiration,
	}
}

// Hash 缓存集合
type Hash struct {
	fields     map[string]any
	expiration int64
}

// Exist Hash中是否存在field
func (h *Hash) Exist(field string) bool {
	if _, exist := h.fields[field]; exist {
		return true
	}
	return false
}

// HSet 添加Hash中的元素
func (h *Hash) HSet(field string, v any) {
	h.fields[field] = v
}

// HGet 获取Hash中的元素
func (h *Hash) HGet(field string) (any, error) {
	if !h.Exist(field) {
		return nil, ErrHashField
	}
	return h.fields[field], nil
}

// HDel 删除Hash中的元素
func (h *Hash) HDel(field string) {
	if h.Exist(field) {
		delete(h.fields, field)
	}
}

// HKeys 获取Hash中所有field列表
func (h *Hash) HKeys() ([]string, error) {
	var fields []string
	for field, _ := range h.fields {
		fields = append(fields, field)
	}
	return fields, nil
}

// HVals 获取Hash中所有的内容列表
func (h *Hash) HVals() ([]any, error) {
	var vals []any
	for _, val := range h.fields {
		vals = append(vals, val)
	}
	return vals, nil
}
