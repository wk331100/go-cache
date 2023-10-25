package go_cache

import (
	"fmt"
	"go-cache/types"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	name1 = "zhangSan"
	name2 = "lisi"
	c     = NewCache()
)

func TestGetSet(t *testing.T) {
	c.Set("name", name1)
	name, err := c.Get("name")
	require.Nil(t, err)
	require.Equal(t, name1, name)
}

func TestIncrDecr(t *testing.T) {
	key := "count"
	c.Incr(key)
	num, err := c.Get(key)
	require.Nil(t, err)
	require.Equal(t, int64(1), num)
	c.Decr(key)
	num, err = c.Get(key)
	require.Nil(t, err)
	require.Equal(t, int64(0), num)
}

func TestIncrByDecrBy(t *testing.T) {
	key := "count"
	c.IncrBy(key, 100)
	num, err := c.Get(key)
	require.Nil(t, err)
	require.Equal(t, int64(100), num)
	c.DecrBy(key, 50)
	num, err = c.Get(key)
	require.Nil(t, err)
	require.Equal(t, int64(50), num)
}

func TestLPushRPopLLen(t *testing.T) {
	key := "queue"
	c.LPush(key, 5)
	l := c.LLen(key)
	require.Equal(t, 1, l)
	c.LPush(key, 4)
	c.LPush(key, 3)
	c.LPush(key, 2)
	c.LPush(key, 1)

	l = c.LLen(key)
	require.Equal(t, 5, l)
	v, err := c.RPop(key)
	require.Nil(t, err)
	require.Equal(t, 5, v)
	v, err = c.RPop(key)
	require.Nil(t, err)
	require.Equal(t, 4, v)
	v, err = c.RPop(key)
	require.Nil(t, err)
	require.Equal(t, 3, v)
	v, err = c.RPop(key)
	require.Nil(t, err)
	require.Equal(t, 2, v)
	v, err = c.RPop(key)
	require.Nil(t, err)
	require.Equal(t, 1, v)
}

func TestLRange(t *testing.T) {
	key := "queue"
	c.LPush(key, 5)
	c.LPush(key, 4)
	c.LPush(key, 3)
	c.LPush(key, 2)
	c.LPush(key, 1)
	l, err := c.LRange(key, 0, 4)
	require.Nil(t, err)
	for index, v := range l {
		require.Equal(t, index+1, v)
	}
}

func TestHGetHSet(t *testing.T) {
	key := "hKey"
	c.HSet(key, "name", name1)
	c.HSet(key, "age", 18)
	name, err := c.HGet(key, "name")
	require.Nil(t, err)
	require.Equal(t, name1, name)
	age, err := c.HGet(key, "age")
	require.Nil(t, err)
	require.Equal(t, 18, age)
}

func TestHKeysValsDel(t *testing.T) {
	key := "hKey"
	c.HSet(key, "name", name1)
	c.HSet(key, "age", 18)
	keys, err := c.HKeys(key)
	require.Nil(t, err)
	require.Equal(t, []string{"name", "age"}, keys)
	vals, err := c.HVals(key)
	require.Nil(t, err)
	require.Equal(t, []any{name1, 18}, vals)
	c.HDel(key, "age")
	keys, err = c.HKeys(key)
	require.Nil(t, err)
	require.Equal(t, []string{"name"}, keys)
}

func TestSAddSRem(t *testing.T) {
	key := "class1"
	m1 := "zhangSan"
	m2 := "liSi"
	c.SAdd(key, m1)
	c.SAdd(key, m2)
	members := c.SCard(key)
	require.Equal(t, 2, members)
	c.SRem(key, m1)
	members = c.SCard(key)
	require.Equal(t, 1, members)
	r1, err := c.SIsMember(key, m1)
	require.Nil(t, err)
	require.False(t, r1)
	r2, err := c.SIsMember(key, m2)
	require.Nil(t, err)
	require.True(t, r2)
}

func TestSUnionSInter(t *testing.T) {
	key1 := "class1"
	key2 := "class2"
	m1 := "zhangSan"
	m2 := "liSi"
	m3 := "wangWu"
	c.SAdd(key1, m1)
	c.SAdd(key1, m2)
	c.SAdd(key2, m2)
	c.SAdd(key2, m3)
	ms, err := c.SMembers(key1)
	require.Nil(t, err)
	require.Equal(t, []any{m1, m2}, ms)
	union := c.SUnion(key1, key2)
	um, _ := union.SMembers()
	require.Equal(t, 3, len(um))
	inter := c.SInter(key1, key2)
	im, _ := inter.SMembers()
	require.Equal(t, 1, len(im))
}

func TestZAddZRem(t *testing.T) {
	key := "english"
	e1 := "zhangSan"
	e2 := "liSi"
	e3 := "wangWu"
	c.ZAdd(key, e1, 100)
	c.ZAdd(key, e2, 90)
	c.ZAdd(key, e3, 95)
	num := c.ZCard(key)
	require.Equal(t, 3, num)
	c.ZRem(key, e3)
	num = c.ZCard(key)
	require.Equal(t, 2, num)
}

func TestZIncrDecr(t *testing.T) {
	key := "english"
	e1 := "zhangSan"
	e2 := "liSi"
	c.ZAdd(key, e1, 100)
	c.ZAdd(key, e2, 90)

	res1 := c.ZIncrBy(key, e1, 20)
	require.Equal(t, float64(120), res1)
	res2 := c.ZDecrBy(key, e2, 10)
	require.Equal(t, float64(80), res2)
}

func TestRank(t *testing.T) {
	key := "rank"
	e1 := "zhangSan"
	e2 := "liSi"
	e3 := "wangWu"
	c.ZAdd(key, e1, 100)
	c.ZAdd(key, e2, 90)
	r1 := c.ZRank(key, e1)
	require.Equal(t, 1, r1)
	r2 := c.ZRank(key, e2)
	require.Equal(t, 2, r2)
	c.ZAdd(key, e3, 95)
	r2 = c.ZRank(key, e2)
	require.Equal(t, 3, r2)
	r1, score := c.ZRankWithScore(key, e1)
	require.Equal(t, 1, r1)
	require.Equal(t, float64(100), score)

	r1 = c.ZRevRank(key, e1)
	require.Equal(t, 3, r1)
	r1, score = c.ZRevRankWithScore(key, e1)
	require.Equal(t, 3, r1)
	require.Equal(t, float64(100), score)
}

func TestRange(t *testing.T) {
	key := "range"
	e1 := "zhangSan"
	e2 := "liSi"
	e3 := "wangWu"
	c.ZAdd(key, e1, 100)
	c.ZAdd(key, e2, 90)
	c.ZAdd(key, e3, 95)
	elements, err := c.ZRange(key, 0, 1)
	require.Nil(t, err)
	require.Equal(t, []string{e1, e3}, elements)
	elements, err = c.ZRevRange(key, 0, 1)
	require.Nil(t, err)
	require.Equal(t, []string{e2, e3}, elements)
	m, err := c.ZRangeWithScore(key, 1, 1)
	require.Nil(t, err)
	require.Equal(t, map[string]float64{e3: 95}, m)
	m, err = c.ZRevRangeWithScore(key, 1, 1)
	require.Nil(t, err)
	require.Equal(t, map[string]float64{e3: 95}, m)
}

func TestExpiration(t *testing.T) {
	k := "exp"
	v := "hello"
	c.SetEx(k, v, time.Microsecond*500)
	v1, err := c.Get(k)
	require.Nil(t, err)
	require.Equal(t, v, v1)
	time.Sleep(time.Second)
	v2, err := c.Get(k)
	require.Nil(t, nil, v2)
	require.Equal(t, types.ErrKeyNotExist, err)
}

// ========== test benchmark ==============

func BenchmarkSetString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	key := "benchmark"

	for j := 0; j < b.N; j++ {
		c.Set(key, "h")
	}
}

func BenchmarkGetString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	key := "benchmark"
	for j := 0; j < b.N; j++ {
		_, _ = c.Get(key)
	}
}

func TestSetStringQPS(t *testing.T) {
	start := time.Now()
	loop := 1000000
	for j := 0; j < loop; j++ {
		c.Set(strconv.Itoa(j), "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set string: %v s\n", d.Seconds())
	fmt.Printf("benchmark set string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}

func TestSetSameStringQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		c.Set(key, "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set same string: %v s\n", d.Seconds())
	fmt.Printf("benchmark set same string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}

func TestGetStringQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		_, _ = c.Get(key)
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark get string: %v s\n", d.Seconds())
	fmt.Printf("benchmark get string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}

func TestSetHashQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		c.HSet(key, strconv.Itoa(j), "h")
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark set hash: %v s\n", d.Seconds())
	fmt.Printf("benchmark set hash qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}

func TestGetHashQPS(t *testing.T) {
	start := time.Now()
	key := "benchmark"
	loop := 1000000
	for j := 0; j < loop; j++ {
		_, _ = c.HGet(key, strconv.Itoa(j))
	}
	d := time.Now().Sub(start)
	fmt.Printf("benchmark get string: %v\n", d.Seconds())
	fmt.Printf("benchmark get string qps: %.0f \n", math.Round(float64(loop)/d.Seconds()))
}
