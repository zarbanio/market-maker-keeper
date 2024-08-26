require("@nomicfoundation/hardhat-toolbox");

const PRIVATE_KEY = process.env.PRIVATE_KEY;
const ARBISCAN_KEY = process.env.ARBISCAN_KEY;

module.exports = {
  solidity: {
    version: "0.8.17",
    settings: {
      optimizer: { enabled: true, runs: 200 },
    },
  },
  etherscan: {
    apiKey: {
      arb: ARBISCAN_KEY,
    },
    customChains: [
      {
        network: "arb",
        chainId: 42161,
        urls: {
          apiURL: "https://api.arbiscan.io/api",
          browserURL: "https://arbiscan.io",
        },
      },
    ],
  },
  networks: {
    arb: {
      url: "https://1rpc.io/arb",
      chainId: 42161,
      accounts: [PRIVATE_KEY],
    },
  },
};
