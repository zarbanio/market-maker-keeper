import {loadFixture} from "@nomicfoundation/hardhat-network-helpers";
import {expect} from "chai";
import {ethers} from "hardhat";
import {FeeAmount, TICK_SPACINGS} from "./shared/constants";
import {getMaxTick, getMinTick} from "./shared/ticks";

const FACTORY_ABI = require("@uniswap/v3-core/artifacts/contracts/UniswapV3Factory.sol/UniswapV3Factory.json").abi;
const SWAP_ROUTER_API = require("@uniswap/v3-periphery/artifacts/contracts/SwapRouter.sol/SwapRouter.json").abi;
const NonfungiblePositionManager_ABI = require("@uniswap/v3-periphery/artifacts/contracts/NonfungiblePositionManager.sol/NonfungiblePositionManager.json").abi;

describe("DexTrader", function () {
    let owner, executor;

    let factory, swapRouter, nonfungiblePositionManager;
    let DAI, ZAR;

    before(async function () {
        [owner, executor] = await ethers.getSigners();

        factory = await ethers.getContractAt(FACTORY_ABI, "0x1F98431c8aD98523631AE4a59f267346ea31F984");
        swapRouter = await ethers.getContractAt(SWAP_ROUTER_API, "0xE592427A0AEce92De3Edee1F18E0157C05861564")
        nonfungiblePositionManager = await ethers.getContractAt(NonfungiblePositionManager_ABI, "0xC36442b4a4522E871399CD717aBDD847Ab11FE88")

        const ERC20 = await ethers.getContractFactory("MintableERC20");
        DAI = await ERC20.deploy('Dai Stablecoin', 'DAI');
        ZAR = await ERC20.deploy('Zar Stablecoin', 'ZAR');

        const amount = ethers.utils.parseEther("1000000")

        await DAI.mint(owner.address, amount);
        await ZAR.mint(owner.address, amount);

        // approve
        await DAI.approve(nonfungiblePositionManager.address, amount)
        await ZAR.approve(nonfungiblePositionManager.address, amount)
    });

    async function deployDexTrader() {
        const DexTrader = await ethers.getContractFactory("DexTrader");
        const dexTrader = await DexTrader.deploy(swapRouter.address, executor.address);

        return dexTrader;
    }

    describe("Deployment", function () {
        it("Should set executor correctly", async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            expect(await dexTrader.executor()).to.equal(executor.address);
        });

        it("Should set swapRouter correctly", async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            expect(await dexTrader.swapRouter()).to.equal(swapRouter.address);
        });

        it("Should set owner correctly", async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            expect(await dexTrader.owner()).to.equal(owner.address);
        });

        it('Should update executor correctly', async function () {
            const dexTrader = await loadFixture(deployDexTrader);
            expect(await dexTrader.executor()).to.equal(executor.address);

            const newExecutor = owner.address;
            await dexTrader.updateExecutor(newExecutor);
            expect(await dexTrader.executor()).to.equal(owner.address);
        });

        it('Should update swapRouter correctly', async function () {
            const dexTrader = await loadFixture(deployDexTrader);
            const newRouter = dexTrader.address;
            await dexTrader.updateRouter(newRouter);
            expect(await dexTrader.swapRouter()).to.equal(newRouter);
        });
    });

    describe("Withdrawals", function () {
        it("Should withdraw erc20 tokens", async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            await DAI.transfer(dexTrader.address, 77);
            expect(await DAI.balanceOf(dexTrader.address)).to.equal(77);
            const balanceBefore = await DAI.balanceOf(owner.address);
            await dexTrader.withdrawToken(DAI.address, 77, owner.address);
            const balanceAfter = await DAI.balanceOf(owner.address);
            expect(balanceAfter.sub(balanceBefore)).to.equal(77);
        });
    });

    describe("Create Pool and Position and Trade", function () {
        let pool
        let token0, token1

        it("", async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            /*********  Should Create Pool  *******/
            console.log('\t* Should create pool')
            // create pool
            if (DAI.address < ZAR.address) {
                token0 = DAI.address
                token1 = ZAR.address
                console.log('\tDai is Token0')
            } else {
                token0 = ZAR.address
                token1 = DAI.address
                console.log('\tDai is Token1')
            }

            // Mock SqrtPriceX96 for create pool
            const SqrtPriceX96 = "594565596886140445191711314693"

            await nonfungiblePositionManager.createAndInitializePoolIfNecessary(token0, token1, FeeAmount.HIGH, SqrtPriceX96);

            pool = await factory.getPool(DAI.address, ZAR.address, FeeAmount.HIGH)

            // check pool address be not Zero
            expect(ethers.utils.isAddress(pool)).to.be.true;
            expect(pool).to.not.equal(ethers.constants.AddressZero);


            /*********  Should Create Position  *******/
            console.log('\t* Should create position');
            const blockNumber = await ethers.provider.getBlockNumber();
            const block = await ethers.provider.getBlock(blockNumber);
            const timestamp = block.timestamp;

            await nonfungiblePositionManager.mint([
                token0,
                token1,
                FeeAmount.HIGH,
                getMinTick(TICK_SPACINGS[FeeAmount.HIGH]),
                getMaxTick(TICK_SPACINGS[FeeAmount.HIGH]),
                ethers.utils.parseEther("1000"), // amount0
                ethers.utils.parseEther("1000"), // amount1
                0, // amount0Min: for check Price slippage
                0, // amount1Min: for check Price slippage
                owner.address,
                timestamp + 10000,
            ]);



            /*********  Should Trade Tokens Correctly  *******/
            console.log('\t* Should trade tokens correctly');

            await DAI.transfer(dexTrader.address, ethers.utils.parseEther("1000"));

            await dexTrader.connect(executor).trade(DAI.address, ZAR.address, FeeAmount.HIGH, ethers.utils.parseEther("1"), 0);

            let zarBalance =  await ZAR.balanceOf(dexTrader.address)
            console.log(`\tDexTrader balance after trade: ${zarBalance}ZAR`)
            await expect(zarBalance).greaterThan(0);
        });

        it('should only executor trade', async function () {
            const dexTrader = await loadFixture(deployDexTrader);

            await DAI.transfer(dexTrader.address, 1000);

            await expect(dexTrader.trade(DAI.address, ZAR.address, FeeAmount.HIGH, 1, 0))
                .to.be.revertedWith("Trader: caller is not the executor")
        });
    });
});
