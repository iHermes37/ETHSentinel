// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Test.sol";

contract Hello {
    function sayHello() public pure returns (string memory) {
        return "Hello World";
    }
}

contract HelloTest is Test {
    function testSayHello() public {
        Hello hello = new Hello();
        string memory greeting = hello.sayHello();

        // 打印到控制台
        emit log_string(greeting);

        // 断言
        assertEq(greeting, "Hello World");
    }
}
