
# Market Maker Keeper

Market Maker Keeper is a bot that helps to keep the Zar stable coin price stable by exploiting arbitrage opportunities between the Zar price on different exchanges. Everyone can participate in the Zar ecosystem by running the Market Maker Keeper bot. The bot is designed to be run on a server and it will automatically buy and sell Zar on different exchanges to keep the price stable.

## Installation

### Build Dependencies:

- [Go](https://golang.org/dl/) (v1.19 or higher)
- [Hardhat](https://hardhat.org/getting-started/) (v2.11.2 or higher)

### abigen:
Install abigen by running the following command. This is required to generate the contract bindings.
```
git clone git@github.com:ethereum/go-ethereum.git
cd go-ethereum
make devtools
```

### Clone the repository:
Clone the repository and update the submodules. This will clone the go-ethereum arbitrum module which is required to build the project.
```
git clone git@github.com:zarbanio/market-maker-keeper.git
cd market-maker-keeper
git submodule update --init --recursive
```

### Code generate: 
```
make code-gen
```

### Install go modules:
```
go mod download
go mod tidy
```

### Compile project:
```
go build main.go
```

### Compile contracts:
Go to the contracts directory and run the following command to compile the contracts.

```shell script
cd dextrader && npm i && npx hardhat compile
```
```shell script
npx hardhat node --fork https://1rpc.io/arb
npx hardhat test --network localhost
```

### Deploy contracts:
```
npx hardhat run scripts/deploy.ts --network main
```

### Configuration:
```
cp configs/local/config.sample.yaml configs/local/config.yaml
```
Edit `configs/local/config.yaml` and set the necessary values.


### Run Locally:
```
go run main.go run --config=configs/local/config.yaml
```

## Contributing

Contributions are welcome! If you have any enhancements, bug fixes, or new features to propose, please submit a pull request.