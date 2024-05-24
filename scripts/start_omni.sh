#!/bin/sh

rm -rf $HOME/.omni/

cd $HOME

chain_id=omnitestnet-1

omnid init --chain-id=$chain_id omni-node --home=$HOME/.omni
echo "bottom soccer blue sniff use improve rough use amateur senior transfer quarter" | omnid keys add validator1 --keyring-backend test --recover
echo "wreck layer draw very fame person frown essence approve lyrics sustain spoon" | omnid keys add validator2 --keyring-backend test --recover
echo "exotic merit wrestle sad bundle age purity ability collect immense place tone" | omnid keys add validator3 --keyring-backend test --recover
omnid add-genesis-account $(omnid keys show validator1 -a --keyring-backend test) 110000000000000uomni
omnid add-genesis-account $(omnid keys show validator2 -a --keyring-backend test) 120000000000000uomni
omnid add-genesis-account $(omnid keys show validator3 -a --keyring-backend test) 130000000000000uomni
omnid gentx validator1 500000000uomni --keyring-backend=test --home=$HOME/.omni --chain-id=$chain_id
omnid collect-gentxs --home=$HOME/.omni
sed -i 's/stake/uomni/g' $HOME/.omni/config/genesis.json
sed -i 's/cors_allowed_origins\s*=\s*\[\]/cors_allowed_origins = ["*",]/g' $HOME/.omni/config/config.toml
sed -i 's/127.0.0.1:26657/0.0.0.0:26657/g' $HOME/.omni/config/config.toml
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME/.omni/config/app.toml
sed -i 's/127.0.0.1:1317/0.0.0.0:1317/g' $HOME/.omni/config/app.toml
sed -i 's/enable = false/enable = true/g' $HOME/.omni/config/app.toml

echo "{
    \"chain_id\": \"omnitestnet-1\",
    \"chain_host\": \"127.0.0.1:1317\",
    \"chain_rpc\": \"127.0.0.1:26657\",
    \"signer_name\": \"validator1\",
    \"signer_passwd\": \"password\",
    \"node_rpc\":\"https://eth-sepolia.g.alchemy.com/v2/yUTFuyJiwaZJmAAFHuoDeEWu-P9tQTSu\"
}" > $HOME/.omni/config.json
omnid start --home=$HOME/.omni