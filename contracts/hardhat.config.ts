import {HardhatUserConfig} from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import {NETWORKS_RPC_URL} from './helper-hardhat-config';

require('dotenv').config()

const Private = process.env.Private || '';
const ARBISCAN_KEY = process.env.ARBISCAN_KEY || '';

const config: HardhatUserConfig = {
    solidity: {
        version: '0.8.17',
        settings: {
            optimizer: {enabled: true, runs: 200},
        },
    },
    typechain: {
        outDir: 'types',
        target: 'ethers-v5',
    },
    etherscan: {
        apiKey: {
            main: ARBISCAN_KEY,
            goerli: ARBISCAN_KEY,
        },
        customChains: [
            {
                network: 'main',
                chainId: 42161,
                urls: {
                    apiURL: 'https://api.arbiscan.io/api',
                    browserURL: 'https://arbiscan.io',
                },
            },
        ]
    },
    networks: {
        main: {
            url: NETWORKS_RPC_URL['main'],
            chainId: 42161,
            accounts: [Private]
        },
    },
};

export default config;
