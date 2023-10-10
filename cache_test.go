package go_cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	name1 = "zhangSan"
	name2 = "lisi"
	c     = NewCache()
)

func TestGetSet(t *testing.T) {
	err := c.Set("name", name1)
	require.Nil(t, err)
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
