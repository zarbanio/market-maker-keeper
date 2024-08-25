require('dotenv').config();

const ALCHEMY_KEY = process.env.ALCHEMY_KEY || '';

export const NETWORKS_RPC_URL = {
    ['main']: ALCHEMY_KEY
        ? `https://arb-mainnet.g.alchemy.com/v2/${ALCHEMY_KEY}`
        : `https://arb1.arbitrum.io/rpc`,
};
