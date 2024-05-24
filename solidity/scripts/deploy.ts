// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// When running the script with `npx hardhat run <script>` you'll find the Hardhat
// Runtime Environment's members available in the global scope.
import {spawn} from "child_process";
import {ethers} from "hardhat";

const deploy = async (name: string, ...args: any) => {
  const Registry = await ethers.getContractFactory(name);
  const registry = await Registry.deploy(...args);
  await registry.deployed();
  console.log(`${name} deployed to:`, registry.address);
  return registry;
};

async function deployProc() {
  /** Deploy contracts **/
  await deploy("Balance");
}

deployProc();
