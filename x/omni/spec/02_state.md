<!--
order: 2
-->

# State

## Params

```proto
message Params {
  option (gogoproto.goproto_stringer) = false;
  // number of block epoch in refreshing balance
	uint64 num_epochs = 1;
  // minimum number of validators voted
	uint64 min_consensus = 2;
  // currnet consensus round
  uint64 cur_round = 3;
  // if whitelisting enabled
  bool is_whitelist_enabled = 4;
}
```

## Balance

```proto
message Balance {
  uint64 index = 1; 
  string balance = 2[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ];
  google.protobuf.Timestamp lastUpdated = 3
    [(gogoproto.stdtime) = true, (gogoproto.nullable) = false]; 
}
```

`Balance` describes the saved balance data of the Ethereum state.

## Observe vote

```proto
message ObserveVote {
  uint64 index = 1; 
  string voter = 2; 
  uint64 round = 3; 
  string value = 4[
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Int",
    (gogoproto.nullable) = false
  ]; 
  google.protobuf.Timestamp timestamp = 5
  [(gogoproto.stdtime) = true, (gogoproto.nullable) = false]; 
}

```

`ObserveVote` describes the observation vote on the Ethereum state storage.
