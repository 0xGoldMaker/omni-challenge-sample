package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// WhitelistKeyPrefix is the prefix to retrieve all Whitelist
	WhitelistKeyPrefix = "Whitelist/value/"
)

// WhitelistKey returns the store key to retrieve a Whitelist from the index fields
func WhitelistKey(
	index uint64,
) []byte {
	var key []byte

	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexBytes, index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
