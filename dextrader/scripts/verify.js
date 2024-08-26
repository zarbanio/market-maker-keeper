const hre = require("hardhat");

const fatalErrors = [
    "The address provided as argument contains a contract, but its bytecode",
    "Daily limit of 100 source code submissions reached",
    "has no bytecode. Is the contract deployed to this network",
    "The constructor for",
];

const okErrors = ["Contract source code already verified"];

const unableVerifyError = 'Fail - Unable to verify';

function delay(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

const verifyArbiscanContract = async (contractAddress, params) => {
    try {
        console.log('Verifying deployed contract');
        const msDelay = 3000;
        const times = 4;

        await verifyWithRetry(contractAddress, params, times, msDelay);
    } catch (error) {
        console.error('[ERROR] Failed to verify contract:', error.message);
    }
}

const verifyWithRetry = async (contractAddress, params, times, msDelay) => {
    let counter = times;
    await delay(msDelay);

    try {
        if (times > 1) {
            await verify(contractAddress, params);
        } else if (times === 1) {
            console.log('Trying to verify via uploading all sources.');
            await verify(contractAddress, params);
        } else {
            console.error('[ERROR] Errors after all the retries, check the logs for more information.');
        }
    } catch (error) {
        counter--;

        if (okErrors.some((okReason) => error.message.includes(okReason))) {
            console.info('Skipping due OK response: ', error.message);
            return;
        }

        if (fatalErrors.some((fatalError) => error.message.includes(fatalError))) {
            console.error('[ERROR] Fatal error detected, skip retries and resume deployment.', error.message);
            return;
        }

        console.error('[ERROR]', error.message);
        console.log();
        console.info(`[INFO] Retrying attempts: ${counter}.`);
        if (error.message.includes(unableVerifyError)) {
            console.log('Trying to verify via uploading all sources.');
            params.relatedSources = undefined;
        }
        await verifyWithRetry(contractAddress, params, counter, msDelay);
    }
}

const verify = async (address, params) => {
    return hre.run("verify:verify", {
        address: address,
        constructorArguments: params,
    });
}

module.exports = {
    verifyArbiscanContract,
};
