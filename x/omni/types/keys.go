package types

const (

	// ModuleName defines the module name
	ModuleName = "omni"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_omni"

	// ParamsKey is the prefix for parameters of omni module
	ParamsKey = "omni_params"

	NativeCoin = "uomni"

	MIN_CONSENSUS = 1
)

var (
	GlobalStoreKeyPrefix = []byte{0x00}
	// WhitelistedNodeCountStoreKey count of whitelisted node store key
	WhitelistedNodeCountStoreKey = append(GlobalStoreKeyPrefix, []byte("WhitelistedNodeCount")...)
	// Observe vote count store key
	ObserveTxCountStoreKey = append(GlobalStoreKeyPrefix, []byte("ObserveTxCount")...)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
