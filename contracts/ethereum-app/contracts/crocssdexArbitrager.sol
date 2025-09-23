// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;
import {IPoolAddressesProvider, IPool} from "@aave/core-v3/contracts/interfaces/IPool.sol";
import {IFlashLoanSimpleReceiver} from "@aave/core-v3/contracts/flashloan/interfaces/IFlashLoanSimpleReceiver.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { IUniswapV2Callee } from "@uniswap/v2-core/contracts/interfaces/IUniswapV2Callee.sol";
import { IUniswapV2Pair } from "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";


contract crocssdex is IFlashLoanSimpleReceiver, IUniswapV2Callee{

    //自定义字段
    IPool private pool;


    constructor(address _provider){
        IPoolAddressesProvider provider=IPoolAddressesProvider(_provider);
        pool=IPool(provider.getPool());
    }

    //------------------外部发起闪电贷和套利程序-------------------------
    function requestFlashLoan(address asset, uint256 amount) external {
        pool.flashLoanSimple(
            address(this), // 接收方
            asset, // 借贷资产
            amount, // 借贷数量
            "", 
            0);
    }

    //-----------必须闪电贷实现的接口-----------------------
    function ADDRESSES_PROVIDER() external view returns (IPoolAddressesProvider) { return IPoolAddressesProvider(address(pool)); } 
    function POOL() external view returns (IPool) { return pool; }
    /** 闪电贷回调 */

     function executeOperation(  
        address asset,//借贷资产的合约地址，例如 USDT 或 ETH 的 ERC20 地址。
        uint256 amount,//本次闪电贷借到的数量。
        uint256 premium,//闪电贷手续费，即你必须在交易结束时归还的额外金额。通常是 amount * feeRate。
        address, // 发起闪电贷的账户地址。通常是调用 flashLoanSimple 的外部地址（你的合约或外部用户）。
        bytes calldata // 可选参数，用于传递自定义数据。例如你想在回调里告诉套利逻辑要操作哪条交易对。
    ) external override returns (bool) {
        // 确认闪电贷到账
        require(IERC20(asset).balanceOf(address(this)) >= amount, "Flashloan failed");


        //跨dex套利程序
        IUniswapV2Pair(pairA).swap(
            pullTokenIs0 ? x : 0,
            pullTokenIs0 ? 0 : x,
            address(this),
            abi.encode(pairB, pullToken, remainToken, y, z)
        );{
            (bool success, ) =
                remainToken.call(abi.encodeWithSignature("transfer(address,uint256)", treasury, y - z - 1));
            require(success, "erc20 transfer 3 failing");
        }

        //归还闪电贷
        IERC20(asset).approve(address(pool), amount + premium);

        return true;
    }

        function uniswapV2Call(
            address sender,
            uint256 amount0,
            uint256 amount1,
            bytes calldata data
        ) external override {
            bool pullTokenIs0 = amount0 > 0;
            uint256 x = pullTokenIs0 ? amount0 : amount1;

            (address pairB, address pullToken, address remainToken, uint256 y, uint256 z) =
                abi.decode(data, (address, address, address, uint256, uint256));

            {
                // IERC20(pullToken).transfer(pairB, x);
                (bool success, ) = pullToken.call(abi.encodeWithSignature("transfer(address,uint256)", pairB, x));
                require(success, "erc20 transfer 1 failing");
            }
            IUniswapV2Pair(pairB).swap(pullTokenIs0 ? 0 : y, pullTokenIs0 ? y : 0, address(this), "");

            {
                // IERC20(remainToken).transfer(pairA, z + 1);
                (bool success, ) =
                    remainToken.call(abi.encodeWithSignature("transfer(address,uint256)", msg.sender, z + 1));
                require(success, "erc20 transfer 2 failing");
            }
        }
}
