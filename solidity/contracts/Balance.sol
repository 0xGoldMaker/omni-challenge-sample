// SPDX-License-Identifier: MIT LICENSE

pragma solidity ^0.8.0;

contract Balance {
    // balance variable
    uint256 private balance;

    // Constructor
    constructor() {
        balance = 100;
    }
    
    // Set balance
    function setBalance(uint256 amount) external {
        balance = amount;
    }

    // Get balance
    function getBalance() external view returns (uint256) {
        return balance;
    }
}