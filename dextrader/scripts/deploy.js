const { ethers } = require("hardhat");
const { verifyArbiscanContract } = require("./verify");

require('dotenv').config()

const UniswapV3Router = '0xE592427A0AEce92De3Edee1F18E0157C05861564';

async function main() {
    const [executor] = await ethers.getSigners();
    console.log(`Deploying contracts with the account: ${executor.address}`);
    console.log(`Account balance: ${ethers.utils.formatEther(await executor.getBalance())} ETH`);

    const DexTrader = await ethers.getContractFactory("DexTrader");
    const dexTrader = await DexTrader.deploy(UniswapV3Router, executor.address);

    await dexTrader.deployed();
    console.log(`DexTrader contract deployed to ${dexTrader.address}`);

    await verifyArbiscanContract(dexTrader.address, [UniswapV3Router, executor])
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
