const { ethers } = require("hardhat");
const qrcode = require('qrcode-terminal'); // Import the qrcode-terminal package

require('dotenv').config();

async function main() {
    const [executor] = await ethers.getSigners();
    console.log(`Deploying contracts with the account: ${executor.address}`);
    console.log(`Account balance: ${ethers.utils.formatEther(await executor.getBalance())} ETH`);

    // Generate and display the QR code for the address
    qrcode.generate(executor.address, { small: true }, (qrCode) => {
        console.log('\nAccount Address QR Code:\n');
        console.log(qrCode);
    });

    console.log('\nCaution!');
    console.log('To fund this account, transfer ETH to the following address on the Arbitrum network:');
    console.log(`\n${executor.address}\n`);
    console.log('Scan the QR code above with your wallet to send ETH.');
    console.log('Ensure you are connected to the Arbitrum network when making the transfer.');
    console.log('Do not send ETH to this address on any other network.');
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});
