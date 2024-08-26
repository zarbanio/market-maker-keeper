# Dex Trader Contract

DexTrader contract is wrapper for UniswapV3 based decentralized exchanges
which makes it easier to interact with UniswapV3 interfaces.


## How to run tests?
For the to pass you need to run them against mainnet-fork or any 
other network of your choice that has UniswapV3 deployed.


```shell script
npx hardhat node --fork https://1rpc.io/arb
npx hardhat test --network localhost
```

### .env template:
```dotenv
# Add wallet private key (for deploy contract)
Private=

# Address of DexTrader contract caller
Executor=

# Address of uniswap v3 router
# mainnet
UniSwapV3Router=0xE592427A0AEce92De3Edee1F18E0157C05861564
#testnet 
#UniSwapV3Router=0xB2413c3F8248DB9085cd3348ECD98d744747Fb3B

# Add Alchemy provider keys
ALCHEMY_KEY=

# Optional Arbiscan key, for automatize the verification of the contracts at Arbiscan
ARBISCAN_KEY=
```

### Deploy Dex Trader

* mainnet
```shell script
npx hardhat run scripts/deploy.ts --network main
```