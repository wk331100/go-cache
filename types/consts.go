package types

const (
	DefaultExpiration = -1
	TypeString        = KeyType("string")
	TypeList          = KeyType("list")
	TypeHash          = KeyType("hash")
	TypeSet           = KeyType("set")
	TypeZSet          = KeyType("zSet")
	DefaultScore      = float64(0)
	ErrorRank         = -1
)

var (
	SupportTypes = []KeyType{
		TypeString, TypeList, TypeHash, TypeSet, TypeZSet,
	}
)

// KeyType 键类型
type KeyType string
