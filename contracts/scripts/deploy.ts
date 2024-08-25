import {ethers} from "hardhat";
import {verifyArbiscanContract} from "./verify";

require('dotenv').config()

const UniSwapV3Router = process.env.UniSwapV3Router || '0xE592427A0AEce92De3Edee1F18E0157C05861564';
const Executor = process.env.Executor || '0x9b4420f91Ae0e4b24D357005A372B08F45Dc9885';

async function main() {
    if (UniSwapV3Router === '') {
        console.error("uniswapV3 router is null");
    }
    if (Executor === '') {
        console.error("executor is null");
    }

    const DexTrader = await ethers.getContractFactory("DexTrader");
    const dexTrader = await DexTrader.deploy(UniSwapV3Router, Executor);

    await dexTrader.deployed();

    console.log(`DexTrader contract deployed to ${dexTrader.address}`);

    // try to verify
    await verifyArbiscanContract(dexTrader.address, [UniSwapV3Router, Executor])
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
