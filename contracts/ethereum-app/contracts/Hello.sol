// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import {IPoolAddressesProvider, IPool} from "@aave/core-v3/contracts/interfaces/IPool.sol";
import {IFlashLoanSimpleReceiver} from "@aave/core-v3/contracts/flashloan/interfaces/IFlashLoanSimpleReceiver.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract FlashLoanExample is IFlashLoanSimpleReceiver {
    IPool private pool;
    address private owner;

    constructor(address _provider) {
        IPoolAddressesProvider provider = IPoolAddressesProvider(_provider);
        pool = IPool(provider.getPool());
        owner = msg.sender;
    }

    /** 发起闪电贷 */
    function requestFlashLoan(address asset, uint256 amount) external {
        pool.flashLoanSimple(address(this), asset, amount, "", 0);
    }

    event FlashLoanExecuted(address asset, uint256 amount, uint256 premium);

    /** 闪电贷回调 */
    function executeOperation(
        address asset,
        uint256 amount,
        uint256 premium,
        address, // initiator
        bytes calldata // params
    ) external override returns (bool) {
        // 简单示例：直接还款
        uint256 totalDebt = amount + premium;

        // 确保合约有足够资产偿还
        IERC20(asset).approve(address(pool), totalDebt);

        emit FlashLoanExecuted(asset, amount, premium);
        return true;
    }

    function ADDRESSES_PROVIDER() external view returns (IPoolAddressesProvider) { return IPoolAddressesProvider(address(pool)); } 
    function POOL() external view returns (IPool) { return pool; }

    function getPool() external view returns (IPool) {
        return pool;
    }

    function getOwner() external view returns (address) {
        return owner;
    }
}
