import { expect } from "chai";
import { Balance } from "../typechain";
import { ethers } from "hardhat";

describe("Balance", function () {
  let balance;
  let owner: { address: string; };
  let balanceContract: Balance;

  beforeEach(async function () {
    // Get the ContractFactory and Signers here.
    balance = await ethers.getContractFactory("Balance");
    [owner] = await ethers.getSigners();

    // To deploy our contract, we just have to call Token.deploy() and await
    // for it to be deployed(), which happens once its transaction has been
    // mined.
    balanceContract = await balance.deploy();
  });

  describe("Balance", function () {
    it ("Initial balance check", async function() {
      // Get the initial balance
      const initialBalance = await balanceContract.getBalance();
      // Compare
      expect(initialBalance).to.equal("100");
    })
    it ("Set balance check", async function() {
      // Set new balance
      await balanceContract.setBalance(1000);
      // Get the balance
      const newBalance = await balanceContract.getBalance();
      // Compare
      expect(newBalance).to.equal("1000");
    })

  });
});
