import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "@nomicfoundation/hardhat-foundry";

import "@nomicfoundation/hardhat-toolbox";

const INFURA_URL = "https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128";

const config: HardhatUserConfig = {
  solidity: "0.8.28",
  networks: {
    hardhat: {
      forking: {
        url: INFURA_URL,
        blockNumber: 23282722,
      },
    },
  },
};

export default config;
