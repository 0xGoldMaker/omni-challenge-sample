# Omni network
**Omni** is a blockchain built using Cosmos SDK and CometBFT and created with [Ignite CLI v0.27.1](https://ignite.com/cli).
It is consists of Cosmos SDK app `omnid` and `observerd`. Observer interacts with Omni as well as Sepolia to observe storage value in a smart contract and broadcast it to Omni. Omni has omni custom module which has several queries and messages to store balance data on external network storage. It also has an option to enable whitelisting which will only accept observations from the whitelisted observers. It also contains parameters such as `epoch` number (refresh rate), smart contract address etc. You can update these parameters using `governance proposals` and so there are couple of `governance handlers` implemented(Update param, Update whitelisted keys etc). There is a testing smart contract included which contains very simple storage read/set functions implemented and you can find it in `solidity` folder. 

## Get started
### Requirements

```
1. Cosmos SDK version v0.47.3 
2. CometBFT v0.37.1
3. Golang v1.20
4. Ignite CLI v0.27.1
```

You can install the prerequisites using the following commands
```
sudo apt update
sudo apt upgrade -y
sudo apt install build-essential jq -y
```

Install golang
```
wget -q -O - https://raw.githubusercontent.com/canha/golang-tools-install-script/master/goinstall.sh | bash -s -- --version 1.20
source ~/.profile
go version
```

### Deploy solidity smart contract
```
1. Add env variables from .env
    * `PRIVATE_KEY` of your wallet
    * `GOERLI_URL` from either alchemy or infuria
2. `npx hardhat run ./scripts/deploy.ts --network sepolia`
3.  Output should look like this
Balance deployed to: 0xcc7F90c440ddBd4B082EE7eAA4e7E82E56869C4B

4. type `npx hardhat verify <ADDRESS> <CONSTRUCTOR_ARGS> --network`
ex: npx hardhat verify 0xcc7F90c440ddBd4B082EE7eAA4e7E82E56869C4B --network sepolia
```

### Configure Observer

Observer needs a configuration file created at the root of Cosmos SDK application folder (`$HOME/.omni`) and the file name is `config.json`. You can use the following bash script or manually fill in the file.
```
echo "{
    \"chain_id\": \"omnitestnet-1\",
    \"chain_host\": \"127.0.0.1:1317\",
    \"chain_rpc\": \"127.0.0.1:26657\",
    \"signer_name\": \"validator1\",
    \"signer_passwd\": \"password\",
    \"node_rpc\":\"[SEPOLIA RPC ENDPOINT]"
}" > $HOME/.omni/config.json
```

### Start Omni node

- **Ignite CLI**

You can run the node using ignite CLI. Please note that you should enable API endpoint and CORS of RPC endpoint of the node.
```
ignite chain serve -v
```
- **Bash script**

There are 2 bash scripts written - single node and multiple nodes on single machine.
You can run a single node using the following commands.

```
cd omni
sudo chmod +x ./scripts/start_omni.sh
./scripts/start_omni.sh
```

You can run the multiple nodes in a single instance using the following commands.
```
cd omni
sudo chmod +x ./scripts/setup_localnet_nodes.sh
sudo chmod +x ./scripts/variables.sh
./scripts/setup_localnet_nodes.sh
```

### Release
There is a make file ready and you can relase the binaries using the `make` commands.
```
make install
```

### Unit test
Unit tests are written to cover all functions inside omni module keeper and types folders.
You can test it using `make` command like the followings.
```
make test-unit
```

### How to test ?
You should run omni node as well as observer in a single machine. Observer will fetch data from a smart contract storage and brodcast it to omni node. Omni node will check observation voted and update its status `balance`.

You can check observation voted using the following commands.
```
omnid q omni list-observe-vote
```

You can check the balance using the following commands.
```
omnid q omni list-balance
```
You can update omni modules parameters using the following governance proposal transaction.
```
omnid tx omni update-params [num epoch] [min consensus] [is whitelist enabled] [smart contract address] --title [title] --summary [summary] --metadata [metadata] --deposit 1000000uomni --chain-id omnitestnet-1 --keyring-backend test --from [wallet_key] -y
```
You can also update whitelisted addresses using the following governance proposal transaction.
```
omnid tx omni whitelisted [address] --title [title] --summary [summary] --metadata [metadata] --deposit 10000000uomni --chain-id omnitestnet-1 --keyring-backend test --from treasury -y
```

You should vote on the proposal to be accepted.
```
omnid tx gov vote [proposal_id] yes --from [wallet_key] --chain-id omnitestnet-1 --keyring-backend test -y
```
### Further development

I would like to make the following improvements.
```
// TOTO:
// 1. Consider weight of each voter
// 2. Consider 2/3+ voter of the total validators
// 3. Apply slahsing logic to give penalty to the validators who hasn't attened the balance voting or who provides fake observerations.
// 4. Make observer to be more robust. At least need to count the block confirmation number in order to avoid block-reorg.
// 5. Add an incentive module to pay rewards to the observers in order to incentive them actively take part in observation. We might give reputation score to each of observer.
```

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
