package types

import (
	"sort"
	"sync"
	"time"
)

// NewZSets 创建Sets类型实例
func NewZSets() *ZSets {
	return &ZSets{
		items: make(map[string]*ZSet),
	}
}

// ZSets 类型数据结构
type ZSets struct {
	mu    sync.RWMutex
	items map[string]*ZSet
}

// Exist 判断k是否存在
func (zs *ZSets) Exist(k string) bool {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	return zs.exist(k)
}

// exist 判断k是否存在
func (zs *ZSets) exist(k string) bool {
	z, exist := zs.items[k]
	if !exist {
		return false
	}
	if z.isExpired() {
		zs.Del(k)
		return false
	}
	return true
}

// ZAdd 向有序集合中添加一个元素
func (zs *ZSets) ZAdd(key, element string, score float64) bool {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	z, exist := zs.items[key]
	if !exist {
		z = newZSet()
	}
	z.ZAdd(element, score)
	zs.items[key] = z
	return exist
}

// ZRem 从有序集合中，删除一个元素
func (zs *ZSets) ZRem(key, element string) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	z, exist := zs.items[key]
	if exist {
		z.ZRem(element)
	}
}

// ZIncrBy 向有序集合中一个元素,增加score
func (zs *ZSets) ZIncrBy(key, element string, score float64) float64 {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	z, exist := zs.items[key]
	if !exist {
		z = newZSet()
	}
	res := z.ZIncrBy(element, score)
	zs.items[key] = z
	return res
}

// ZDecrBy 向有序集合中一个元素,减少score
func (zs *ZSets) ZDecrBy(key, element string, score float64) float64 {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	z, exist := zs.items[key]
	if !exist {
		z = newZSet()
	}
	res := z.ZDecrBy(element, score)
	zs.items[key] = z
	return res
}

// ZCard 获取有序集合的元素数量
func (zs *ZSets) ZCard(key string) int {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	z, exist := zs.items[key]
	if !exist {
		return 0
	} else if z.isExpired() {
		zs.Del(key)
		return 0
	}
	return z.ZCard()
}

// ZRank 获取有序集合的元素排名
func (zs *ZSets) ZRank(key, element string) int {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return ErrorRank
	} else if z.isExpired() {
		zs.Del(key)
		return ErrorRank
	}
	return z.ZRank(element)
}

// ZRankWithScore 获取有序集合的元素排名和score
func (zs *ZSets) ZRankWithScore(key, element string) (int, float64) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return ErrorRank, DefaultScore
	} else if z.isExpired() {
		zs.Del(key)
		return ErrorRank, DefaultScore
	}
	return z.ZRankWithScore(element)
}

// ZRevRank 获取有序集合的元素倒数排名
func (zs *ZSets) ZRevRank(key, element string) int {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return ErrorRank
	} else if z.isExpired() {
		zs.Del(key)
		return ErrorRank
	}
	return z.ZRevRank(element)
}

// ZRevRankWithScore 获取有序集合的元素倒数排名和score
func (zs *ZSets) ZRevRankWithScore(key, element string) (int, float64) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return ErrorRank, DefaultScore
	} else if z.isExpired() {
		zs.Del(key)
		return ErrorRank, DefaultScore
	}
	return z.ZRevRankWithScore(element)
}

// ZRange 获取有序集合区间元素
func (zs *ZSets) ZRange(key string, start, stop int) ([]string, error) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return nil, ErrZSetKey
	} else if z.isExpired() {
		zs.Del(key)
		return nil, ErrZSetKey
	}
	if start > stop {
		return nil, ErrStartStop
	}
	return z.ZRange(start, stop), nil
}

// ZRangeWithScore 获取有序集合区间元素包含Score
func (zs *ZSets) ZRangeWithScore(key string, start, stop int) (map[string]float64, error) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return nil, ErrZSetKey
	} else if z.isExpired() {
		zs.Del(key)
		return nil, ErrZSetKey
	}
	if start > stop {
		return nil, ErrStartStop
	}
	return z.ZRangeWithScores(start, stop), nil
}

// ZRevRange 获取有序集合倒排区间元素
func (zs *ZSets) ZRevRange(key string, start, stop int) ([]string, error) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return nil, ErrZSetKey
	} else if z.isExpired() {
		zs.Del(key)
		return nil, ErrZSetKey
	}
	if start > stop {
		return nil, ErrStartStop
	}
	return z.ZRevRange(start, stop), nil
}

// ZRevRangeWithScore 获取有序集合倒排区间元素包含Score
func (zs *ZSets) ZRevRangeWithScore(key string, start, stop int) (map[string]float64, error) {
	zs.mu.RLock()
	defer zs.mu.RUnlock()
	z, exist := zs.items[key]
	if !exist {
		return nil, ErrZSetKey
	} else if z.isExpired() {
		zs.Del(key)
		return nil, ErrZSetKey
	}
	if start > stop {
		return nil, ErrStartStop
	}
	return z.ZRevRangeWithScore(start, stop), nil
}

// Del 删除一个key
func (zs *ZSets) Del(k string) {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	delete(zs.items, k)
}

// Expiration 设置超时时间
func (zs *ZSets) Expiration(k string, d time.Duration) error {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	if !zs.exist(k) {
		return ErrKeyNotExist
	}
	zs.items[k].expiration = time.Now().Add(d).UnixNano()
	return nil
}

// ClearExpiration 清理过期的key
func (zs *ZSets) ClearExpiration() {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	for key, item := range zs.items {
		if item.isExpired() {
			delete(zs.items, key)
		}
	}
}

// RandomClearExpiration 随机清理过期的key
func (zs *ZSets) RandomClearExpiration() {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	var counter int
	for key, item := range zs.items {
		if counter > DefaultCleanItems {
			return
		}
		if item.isExpired() {
			delete(zs.items, key)
		}
		counter++
	}
}

// Flush 清空缓存
func (zs *ZSets) Flush() {
	zs.mu.Lock()
	defer zs.mu.Unlock()
	zs.items = make(map[string]*ZSet)
}

// NewZSet 创建一个集合的实例
func newZSet() *ZSet {
	return &ZSet{
		elements: make(map[string]float64),
	}
}

// ZSet 缓存集合
type ZSet struct {
	elements   map[string]float64
	sorted     []string
	expiration int64
}

// ZAdd 向有序集合中添加一个元素
func (z *ZSet) ZAdd(e string, score float64) {
	if _, exist := z.elements[e]; exist {
		z.ZRevRank(e)
	}
	z.elements[e] = score
	z.sorted = append(z.sorted, e)
	sort.Slice(z.sorted, func(i, j int) bool {
		return z.elements[z.sorted[i]] > z.elements[z.sorted[j]]
	})
}

// ZRem 从有序集合中，删除一个元素
func (z *ZSet) ZRem(e string) {
	delete(z.elements, e)
	for i, m := range z.sorted {
		if m == e {
			z.sorted = append(z.sorted[:i], z.sorted[i+1:]...)
			break
		}
	}
}

// ZIncrBy 向有序集合中一个元素,增加score
func (z *ZSet) ZIncrBy(e string, score float64) float64 {
	if _, exist := z.elements[e]; !exist {
		z.elements[e] = DefaultScore
	}
	z.elements[e] += score
	return z.elements[e]
}

// ZDecrBy 向有序集合中一个元素,减少score
func (z *ZSet) ZDecrBy(e string, score float64) float64 {
	if _, exist := z.elements[e]; !exist {
		z.elements[e] = DefaultScore
	}
	z.elements[e] -= score
	return z.elements[e]
}

// ZCard 获取有序集合的元素数量
func (z *ZSet) ZCard() int {
	return len(z.sorted)
}

// ZRank 获取有序集合的元素排名
func (z *ZSet) ZRank(e string) int {
	for rank, element := range z.sorted {
		if element == e {
			return rank + 1
		}
	}
	return ErrorRank
}

// ZRankWithScore 获取有序集合的元素排名和score
func (z *ZSet) ZRankWithScore(e string) (int, float64) {
	for rank, element := range z.sorted {
		if e == element {
			return rank + 1, z.elements[e]
		}
	}
	return ErrorRank, DefaultScore
}

// ZRevRank 获取有序集合的元素倒数排名
func (z *ZSet) ZRevRank(e string) int {
	for i := len(z.sorted) - 1; i >= 0; i-- {
		if z.sorted[i] == e {
			return len(z.sorted) - i
		}
	}
	return ErrorRank
}

// ZRevRankWithScore 获取有序集合的元素倒数排名和score
func (z *ZSet) ZRevRankWithScore(e string) (int, float64) {
	for i := len(z.sorted) - 1; i >= 0; i-- {
		if z.sorted[i] == e {
			return len(z.sorted) - i, z.elements[e]
		}
	}
	return ErrorRank, DefaultScore
}

// ZRange 获取有序集合区间元素
func (z *ZSet) ZRange(start, stop int) []string {
	if start < 0 {
		start = len(z.sorted) + start
	}
	if stop < 0 {
		stop = len(z.sorted) + stop
	}
	if start > stop || start >= len(z.sorted) {
		return nil
	}
	if stop >= len(z.sorted) {
		stop = len(z.sorted) - 1
	}
	return z.sorted[start : stop+1]
}

// ZRangeWithScores 获取有序集合区间元素包含Score
func (z *ZSet) ZRangeWithScores(start, stop int) map[string]float64 {
	elements := z.ZRange(start, stop)
	result := make(map[string]float64, len(elements))
	for _, element := range elements {
		result[element] = z.elements[element]
	}
	return result
}

// ZRevRange 获取有序集合倒排区间元素
func (z *ZSet) ZRevRange(start, stop int) []string {
	if start < 0 {
		start = len(z.sorted) + start
	}
	if stop < 0 {
		stop = len(z.sorted) + stop
	}
	if start > stop || start >= len(z.sorted) {
		return nil
	}
	if stop >= len(z.sorted) {
		stop = len(z.sorted) - 1
	}
	// Reverse the slice
	rev := make([]string, len(z.sorted))
	copy(rev, z.sorted)
	for i, j := 0, len(rev)-1; i < j; i, j = i+1, j-1 {
		rev[i], rev[j] = rev[j], rev[i]
	}
	return rev[start : stop+1]
}

// ZRevRangeWithScore 获取有序集合倒排区间元素包含Score
func (z *ZSet) ZRevRangeWithScore(start, stop int) map[string]float64 {
	elements := z.ZRevRange(start, stop)
	result := make(map[string]float64, len(elements))
	for _, element := range elements {
		result[element] = z.elements[element]
	}
	return result
}

// isExpired 判断一个元素是否过期
func (z *ZSet) isExpired() bool {
	if time.Now().UnixNano() > z.expiration {
		return true
	}
	return false
}
