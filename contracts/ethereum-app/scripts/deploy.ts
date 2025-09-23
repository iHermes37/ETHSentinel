import { ethers } from "hardhat";

import 'global-agent/bootstrap';
process.env.GLOBAL_AGENT_HTTP_PROXY = 'http://192.168.184.215:7890';

async function main() {
  const [deployer] = await ethers.getSigners();
  console.log("部署者地址:", deployer.address);

  const providerAddress = "0x2f39d218133AFaB8F2B819B1066c7E434Ad94E9e";

  const FlashLoan = await ethers.getContractFactory("FlashLoanExample");
  
  // v6 直接 await deploy 就返回 Contract
  const flashLoan = await FlashLoan.deploy(providerAddress);

  // 不再需要 flashLoan.deployed()
  console.log("FlashLoanExample 部署完成:", flashLoan.target); 
  // target 是合约地址
}

main().catch(console.error);
