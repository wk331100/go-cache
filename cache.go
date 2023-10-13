package go_cache

import (
	"go-cache/types"
	"time"
)

// NewCache 创建新的缓存对象
func NewCache() *Cache {
	return &Cache{
		keyMap:  make(map[string]types.KeyType),
		strings: types.NewStrings(),
		lists:   types.NewLists(),
		hashes:  types.NewHashes(),
		sets:    types.NewSets(),
		zSets:   types.NewZSets(),
	}
}

// Cache 缓存结构
// items 为string类型的映射
type Cache struct {
	keyMap  map[string]types.KeyType
	strings *types.Strings
	lists   *types.Lists
	hashes  *types.Hashes
	sets    *types.Sets
	zSets   *types.ZSets
}

// ======== 字符串 =======

// Set 缓存k的值为v
func (c *Cache) Set(k string, v any) {
	c.strings.Set(k, v)
}

// SetEx 缓存k的值为v,并且设置超时时间d
func (c *Cache) SetEx(k string, v any, d time.Duration) {
	c.strings.SetEx(k, v, d)
}

// Get 获取一个string类型值
func (c *Cache) Get(k string) (any, error) {
	return c.strings.Get(k)
}

// Incr 对k计数+1
func (c *Cache) Incr(k string) {
	c.strings.Incr(k)
}

// Decr 对k计数-1
func (c *Cache) Decr(k string) {
	c.strings.Decr(k)
}

// IncrBy 对k计数+v
func (c *Cache) IncrBy(k string, v int64) {
	c.strings.IncrBy(k, v)
}

// DecrBy 对k计数-v
func (c *Cache) DecrBy(k string, v int64) {
	c.strings.DecrBy(k, v)
}

// ======== 列表 =======

// LPush 从队列k的头部，添加一个元素v
func (c *Cache) LPush(k string, v any) {
	c.lists.LPush(k, v)
}

// LPop 从队列k的头部，弹出一个元素
func (c *Cache) LPop(k string) (any, error) {
	return c.lists.LPop(k)
}

// RPush 从队列k的尾部，添加一个元素
func (c *Cache) RPush(k string, v any) {
	c.lists.RPush(k, v)
}

// RPop 从队列k的尾部，弹出一个元素
func (c *Cache) RPop(k string) (any, error) {
	return c.lists.RPop(k)
}

// LLen 获取队列k的长度
func (c *Cache) LLen(k string) int {
	return c.lists.LLen(k)
}

// LRange 获取队列元素列表
func (c *Cache) LRange(k string, start, stop int) ([]any, error) {
	return c.lists.LRange(k, start, stop)
}

// ======== 散列Hash =======

// HSet 缓存数据到Hash中
// k 为Hash中的key
// field 为hash中项
func (c *Cache) HSet(k, field string, v any) {
	c.hashes.HSet(k, field, v)
}

// HGet 从Hash中获取存储的元素
func (c *Cache) HGet(k, field string) (any, error) {
	return c.hashes.HGet(k, field)
}

// HDel 从Hash中删除元素field
func (c *Cache) HDel(k, field string) {
	c.hashes.HDel(k, field)
}

// HKeys 获取Hash中的所有元素field
func (c *Cache) HKeys(k string) ([]string, error) {
	return c.hashes.HKeys(k)
}

// HVals 获取Hash中所有元素的内容
func (c *Cache) HVals(k string) ([]any, error) {
	return c.hashes.HVals(k)
}

// ======== 集合 =======

// SAdd 向集合中添加一个元素
func (c *Cache) SAdd(k string, m any) {
	c.sets.SAdd(k, m)
}

// SRem 从集合中，删除一个元素
func (c *Cache) SRem(k, m string) {
	c.sets.SRem(k, m)
}

// SMembers 获取集合中所有的元素列表
func (c *Cache) SMembers(k string) ([]any, error) {
	return c.sets.SMembers(k)
}

// SIsMember 判断m是否为集合中的元素
func (c *Cache) SIsMember(k string, m any) (bool, error) {
	return c.sets.SIsMember(k, m)
}

// SCard 统计集合中元素数量
func (c *Cache) SCard(k string) int {
	return c.sets.SCard(k)
}

// SUnion 获取集合s1和s2的并集
func (c *Cache) SUnion(k1, k2 string) *types.Set {
	return c.sets.SUnion(k1, k2)
}

// SInter 获取集合s1和s2的交集
func (c *Cache) SInter(k1, k2 string) *types.Set {
	return c.sets.SInter(k1, k2)
}

// ======== 有序集合 =======

// ZAdd 向有序集合中添加一个元素
func (c *Cache) ZAdd(key, element string, score float64) {
	c.zSets.ZAdd(key, element, score)
}

// ZRem 从有序集合中，删除一个元素
func (c *Cache) ZRem(key, element string) {
	c.zSets.ZRem(key, element)
}

// ZIncrBy 向有序集合中一个元素,增加score
func (c *Cache) ZIncrBy(key, element string, score float64) float64 {
	return c.zSets.ZIncrBy(key, element, score)
}

// ZDecrBy 向有序集合中一个元素,减少score
func (c *Cache) ZDecrBy(key, element string, score float64) float64 {
	return c.zSets.ZDecrBy(key, element, score)
}

// ZCard 获取有序集合的元素数量
func (c *Cache) ZCard(key string) int {
	return c.zSets.ZCard(key)
}

// ZRank 获取有序集合的元素排名
func (c *Cache) ZRank(key, element string) int {
	return c.zSets.ZRank(key, element)
}

// ZRankWithScore 获取有序集合的元素排名和score
func (c *Cache) ZRankWithScore(key, element string) (int, float64) {
	return c.zSets.ZRankWithScore(key, element)
}

// ZRevRank 获取有序集合的元素倒数排名
func (c *Cache) ZRevRank(key, element string) int {
	return c.zSets.ZRevRank(key, element)
}

// ZRevRankWithScore 获取有序集合的元素倒数排名和score
func (c *Cache) ZRevRankWithScore(key, element string) (int, float64) {
	return c.zSets.ZRevRankWithScore(key, element)
}

// ZRange 获取有序集合区间元素
func (c *Cache) ZRange(key string, start, stop int) ([]string, error) {
	return c.zSets.ZRange(key, start, stop)
}

// ZRangeWithScore 获取有序集合区间元素包含Score
func (c *Cache) ZRangeWithScore(key string, start, stop int) (map[string]float64, error) {
	return c.zSets.ZRangeWithScore(key, start, stop)
}

// ZRevRange 获取有序集合倒排区间元素
func (c *Cache) ZRevRange(key string, start, stop int) ([]string, error) {
	return c.zSets.ZRevRange(key, start, stop)
}

// ZRevRangeWithScore 获取有序集合倒排区间元素包含Score
func (c *Cache) ZRevRangeWithScore(key string, start, stop int) (map[string]float64, error) {
	return c.zSets.ZRevRangeWithScore(key, start, stop)
}

// ======== 全局 =======

// Exists 判断key是否存在
func (c *Cache) Exists(k string) bool {
	if _, exist := c.keyMap[k]; exist {
		return true
	}
	return false
}

// HExists 判断Hash中是否存在该field
func (c *Cache) HExists(k, field string) bool {
	if _, err := c.hashes.HGet(k, field); err == nil {
		return true
	}
	return false
}

// Del 删除一个key
func (c *Cache) Del(k string) {
	t := c.keyMap[k]
	switch t {
	case types.TypeString:
		c.strings.Del(k)
	case types.TypeHash:
		c.hashes.Del(k)
	case types.TypeList:
		c.lists.Del(k)
	case types.TypeSet:
		c.sets.Del(k)
	case types.TypeZSet:
		c.zSets.Del(k)
	}
}

// Expiration 设置超时时间
func (c *Cache) Expiration(k string, d time.Duration) error {
	return nil
}
