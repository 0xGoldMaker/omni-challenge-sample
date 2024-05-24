package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// BalanceKeyPrefix is the prefix to retrieve all Balance
	BalanceKeyPrefix = "Balance/value/"
)

// BalanceKey returns the store key to retrieve a Balance from the index fields
func BalanceKey(
	index uint64,
) []byte {
	var key []byte

	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexBytes, index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
