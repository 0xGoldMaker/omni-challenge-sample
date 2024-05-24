# Deployment

1. Add env variables from .env
    * `PRIVATE_KEY` of your wallet
    * `GOERLI_URL` from either alchemy or infuria
2. `npx hardhat run ./scripts/deploy.ts --network sepolia`
3.  Output should look like this
Balance deployed to: 0xcc7F90c440ddBd4B082EE7eAA4e7E82E56869C4B

4. type `npx hardhat verify <ADDRESS> <CONSTRUCTOR_ARGS> --network`
ex: npx hardhat verify 0xcc7F90c440ddBd4B082EE7eAA4e7E82E56869C4B --network sepolia
