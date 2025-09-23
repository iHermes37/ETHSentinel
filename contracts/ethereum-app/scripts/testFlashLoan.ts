import { ethers, network } from "hardhat";
import { FlashLoanExample, IERC20 } from "../typechain-types";

async function main() {
  const [user] = await ethers.getSigners();
  console.log("🚀 测试账户:", user.address);

  // ================= 部署合约 =================
  const providerAddress = "0x2f39d218133AFaB8F2B819B1066c7E434Ad94E9e"; // Aave V3 Mainnet
  const FlashLoan = await ethers.getContractFactory("FlashLoanExample");
  const flashLoan = await FlashLoan.deploy(providerAddress);
  console.log("✅ FlashLoanExample 部署完成:", flashLoan.target);

  // ================= impersonate WETH 大户 =================
  const WETH_ADDRESS = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2";
  const WETH_WHALE = "0x28C6c06298d514Db089934071355E5743bf21d60";
  const WETH = (await ethers.getContractAt("IERC20", WETH_ADDRESS)) as IERC20;

  await network.provider.request({
    method: "hardhat_impersonateAccount",
    params: [WETH_WHALE],
  });
  const whaleSigner = await ethers.getSigner(WETH_WHALE);

  // ================= 给测试账户充值 WETH =================
  const amountToSend = ethers.parseEther("5");
  await WETH.connect(whaleSigner).transfer(user.address, amountToSend);
  console.log(`✅ 转账成功: ${ethers.formatEther(amountToSend)} WETH -> ${user.address}`);

  // ================= 给合约充值 WETH 用于偿还闪电贷 =================
  await WETH.connect(user).transfer(flashLoan.target, ethers.parseEther("2"));
  console.log("💰 已给合约充值 2 WETH 用于偿还闪电贷");

  // ================= 请求闪电贷 =================
  const loanAmount = ethers.parseEther("1"); // 1 WETH
  const tx = await flashLoan.connect(user).requestFlashLoan(WETH_ADDRESS, loanAmount);
  
  // 等待交易被挖矿，获取交易回执
  const receipt = await tx.wait();
  if (!receipt) {
    throw new Error("交易回执为空，闪电贷交易可能失败");
    }

  console.log("📤 闪电贷请求交易发送成功");

  // ================= 捕获闪电贷事件 =================


    const events = await flashLoan.queryFilter(
    flashLoan.filters.FlashLoanExecuted(),
    receipt.blockNumber, // 使用区块号代替 blockHash
    receipt.blockNumber
    );

  events.forEach((e) => {
    console.log("⚡ 闪电贷执行完成:");
    console.log("   asset:", e.args?.asset);
    console.log("   amount:", ethers.formatEther(e.args?.amount));
    console.log("   premium:", ethers.formatEther(e.args?.premium));
  });

  // ================= 验证测试账户 WETH 余额 =================
  const balance = await WETH.balanceOf(user.address);
  console.log("💰 测试账户 WETH 余额:", ethers.formatEther(balance));
}

main()
  .then(() => process.exit(0))
  .catch((err) => {
    console.error(err);
    process.exit(1);
  });
