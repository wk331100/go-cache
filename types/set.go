package types

import (
	"fmt"
	"sync"
	"time"
)

// NewSets 创建Sets类型实例
func NewSets() *Sets {
	return &Sets{
		items: make(map[string]*Set),
	}
}

// Sets 类型数据结构
type Sets struct {
	mu    sync.Mutex
	items map[string]*Set
}

// Exist 判断k是否存在
func (ss *Sets) Exist(k string) bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.exist(k)
}

// exist 判断k是否存在
func (ss *Sets) exist(k string) bool {
	s, exist := ss.items[k]
	if !exist {
		return false
	}
	if s.isExpired() {
		ss.del(k)
		return false
	}
	return true
}

// SAdd 向集合中添加一个元素
func (ss *Sets) SAdd(k string, m any) bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	s, exist := ss.items[k]
	if !exist {
		s = newSet()
	}
	s.SAdd(m)
	ss.items[k] = s
	return exist
}

// SRem 从集合中，删除一个元素
func (ss *Sets) SRem(k, m string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s, exist := ss.items[k]
	if exist {
		s.SRem(m)
	}
}

// SMembers 获取集合中所有的元素列表
func (ss *Sets) SMembers(k string) ([]any, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s, exist := ss.items[k]
	if !exist {
		return nil, ErrSetKey
	} else if s.isExpired() {
		ss.del(k)
		return nil, ErrSetKey
	}
	return s.SMembers()
}

// SIsMember 判断m是否为集合中的元素
func (ss *Sets) SIsMember(k string, m any) (bool, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s, exist := ss.items[k]
	if !exist {
		return false, ErrSetKey
	} else if s.isExpired() {
		ss.del(k)
		return false, ErrSetKey
	}
	return s.SIsMember(m)
}

// SCard 统计集合中元素数量
func (ss *Sets) SCard(k string) int {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s, exist := ss.items[k]
	if !exist {
		return 0
	} else if s.isExpired() {
		fmt.Println("111")
		ss.del(k)
		return 0
	}
	return s.SCard()
}

// SUnion 获取集合s1和s2的并集
func (ss *Sets) SUnion(k1, k2 string) *Set {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s1, exist1 := ss.items[k1]
	if exist1 && s1.isExpired() {
		exist1 = false
	}
	s2, exist2 := ss.items[k2]
	if exist2 && s2.isExpired() {
		exist2 = false
	}
	if !exist1 && !exist2 {
		return nil
	} else if exist1 && !exist2 {
		return s1
	} else if !exist1 && exist2 {
		return s2
	}
	return s1.SUnion(s2)
}

// SInter 获取集合s1和s2的交集
func (ss *Sets) SInter(k1, k2 string) *Set {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	s1, exist1 := ss.items[k1]
	if exist1 && s1.isExpired() {
		exist1 = false
	}
	s2, exist2 := ss.items[k2]
	if exist2 && s2.isExpired() {
		exist2 = false
	}
	if !exist1 || !exist2 {
		return nil
	}
	return s1.SInter(s2)
}

// Del 删除一个key
func (ss *Sets) Del(k string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.del(k)
}

func (ss *Sets) del(k string) {
	delete(ss.items, k)
}

// Expiration 设置超时时间
func (ss *Sets) Expiration(k string, d time.Duration) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.exist(k) {
		return ErrKeyNotExist
	}
	ss.items[k].expiration = time.Now().Add(d).UnixNano()
	return nil
}

// ClearExpiration 清理过期的key
func (ss *Sets) ClearExpiration() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	for key, item := range ss.items {
		if item.isExpired() {
			delete(ss.items, key)
		}
	}
}

// RandomClearExpiration 随机清理过期的key
func (ss *Sets) RandomClearExpiration() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	var counter int
	for key, item := range ss.items {
		if counter > DefaultCleanItems {
			return
		}
		if item.isExpired() {
			delete(ss.items, key)
		}
		counter++
	}
}

// Flush 清空缓存
func (ss *Sets) Flush() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.items = make(map[string]*Set)
}

// newSet 创建一个集合的实例
func newSet() *Set {
	return &Set{
		sets:       make(map[any]struct{}),
		expiration: DefaultExpiration,
	}
}

// Set 缓存集合
type Set struct {
	sets       map[any]struct{}
	expiration int64
}

// SAdd 向集合中添加一个元素
func (s *Set) SAdd(m any) {
	s.sets[m] = struct{}{}
}

// SRem 从集合中，删除一个元素
func (s *Set) SRem(m string) {
	delete(s.sets, m)
}

// SMembers 获取集合中所有的元素列表
func (s *Set) SMembers() ([]any, error) {
	members := make([]any, 0, len(s.sets))
	for m, _ := range s.sets {
		members = append(members, m)
	}
	return members, nil
}

// SIsMember 判断m是否为集合中的元素
func (s *Set) SIsMember(m any) (bool, error) {
	_, isMember := s.sets[m]
	return isMember, nil
}

// SCard 统计集合中元素数量
func (s *Set) SCard() int {
	return len(s.sets)
}

// SUnion 获取集合s1和s2的并集
func (s *Set) SUnion(other *Set) *Set {
	union := newSet()
	for member, _ := range s.sets {
		union.SAdd(member)
	}
	for member, _ := range other.sets {
		union.SAdd(member)
	}
	return union
}

// SInter 获取集合s1和s2的交集
func (s *Set) SInter(other *Set) *Set {
	inter := newSet()
	for member, _ := range other.sets {
		if _, exist := s.sets[member]; exist {
			inter.SAdd(member)
		}
	}
	return inter
}

// isExpired 判断一个元素是否过期
func (s *Set) isExpired() bool {
	if s.expiration != DefaultExpiration && time.Now().UnixNano() > s.expiration {
		return true
	}
	return false
}
