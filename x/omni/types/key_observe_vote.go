package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ObserveVoteKeyPrefix is the prefix to retrieve all ObserveVote
	ObserveVoteKeyPrefix = "ObserveVote/value/"
)

// ObserveVoteKey returns the store key to retrieve a ObserveVote from the index fields
func ObserveVoteKey(
	index uint64,
) []byte {
	var key []byte

	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexBytes, index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
