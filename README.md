<div id="top"></div>


<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">geth-leveldb-explorer</h3>

  <p align="center">
    Database explorer for Geth in GO.
  </p>
</div>


<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* npm
  ```sh
  npm install npm@latest -g
  ```

### Installation TODO

1. Get a free API Key at [https://example.com](https://example.com)
2. Clone the repo
   ```sh
   git clone https://github.com/github_username/repo_name.git
   ```
3. Install NPM packages
   ```sh
   npm install
   ```
4. Enter your API in `config.js`
   ```js
   const API_KEY = 'ENTER YOUR API';
   ```


<!-- USAGE EXAMPLES -->
## Usage

## 1 - LevelDB

## 1.1 - TrieDetails
Search in levelDB the merkle-patricia trees and detail the last one

```sh
  go run main.go trieDetails <LevelDB path>
```

Returns:
 * Total number of state trees (for blocks present in levelDB).
 * Gives the block number and the root of the most recent state tree
 * Total number of accounts (including smartcontract) in the tree
 * Total number of smartcontract in the tree
 * Size of the most recent state tree with leaf details
 * Size of most recent storage tree with leaf details

Example :
```sh
go run main.go trieDetails .ethereum/geth/chaindata/
  
  [...]

  Total number of tree state : 1

  Latest state tree : 
  - Block number : 63e46e
  - State root : 93c3aa9ee4c6285fbe9d28dfbfa245912220dac8fda9c0ecf44ee9677a5f7b19


  Latest state leaf size : 1311662664 bytes
  Latest state tree size : 1856986567 bytes

  Final account number :9302568
  Final smartcontract number :3142527

  Latest storage leaf size : 7044254464 bytes
  Latest storage tree size : 12559033645 bytes
```

## 1.2 - CountStateTrees

Count in levelDB the merkle-patricia trees

```sh
  go run main.go countStateTrees <LevelDB path>
```

Return the total number of state trees (for blocks present in levelDB).

Example :
```sh
  go run main.go countStateTrees .ethereum/geth/chaindata/
  
  [...]
  
  Total number of tree state : 1
```

## 1.3 - SnapshotAccount
Search for an account in the snapshot part of LevelDB

```sh
  go run main.go snapshotAccount <LevelDB path> <account address>
```

Return raw and decoded informations about the account

Example :
```sh
  go run main.go snapshotAccount .ethereum/geth/chaindata/ 8c5fecdC472E27Bc447696F431E425D02dd46a8c
  
  [...]
  
  Snapshot : 
  key : 619a66eb0f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293
  value : f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470

  address : 8c5fecdC472E27Bc447696F431E425D02dd46a8c
  data : {1 c758d0418cfd96ed0000 3651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834 7ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470}
```

## 1.4 - TreeAccount
Search for an account in the merkle-patricia tree part of LevelDB

```sh
  go run main.go treeAccount <LevelDB path> <account address>
```

Return raw and decoded informations about the account

Example :
```sh
  go run main.go treeAccount .ethereum/geth/chaindata/ 8c5fecdC472E27Bc447696F431E425D02dd46a8c
  
  [...]
  
  Merkle-Patricia tree : 
  key : 45afc616075ec2b73fd61a0bd140b7acbda2aca54dd847a610bc4b2cfe4b6ecc
  value : f8709d3f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293b850f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470

  address : 0x8c5fecdC472E27Bc447696F431E425D02dd46a8c
  data : [3f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293 f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470]
  account data : [01 c758d0418cfd96ed0000 3651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834 7ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470]
```
## 1.5 - CompareAccount
Search for an account in the merkle-patricia tree and snapshot in LevelDB

```sh
  go run main.go treeAccount <LevelDB path> <account address>
```

Return raw and decoded informations about the account for both part

Example :
```sh
  go run main.go compareAccount .ethereum/geth/chaindata/ 8c5fecdC472E27Bc447696F431E425D02dd46a8c
  
LevelDB ok
Merkle-Patricia tree : 
key : 45afc616075ec2b73fd61a0bd140b7acbda2aca54dd847a610bc4b2cfe4b6ecc
value : f8709d3f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293b850f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470

address : 0x8c5fecdC472E27Bc447696F431E425D02dd46a8c
data : [3f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293 f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470]
account data : [01 c758d0418cfd96ed0000 3651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834 7ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470]

LevelDB ok
Snapshot : 
key : 619a66eb0f03c4b8bcb5b2c0947f2d835e7f4f740de39eb1e8c38510d275cfa293
value : f84e018ac758d0418cfd96ed0000a03651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834a07ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470

address : 8c5fecdC472E27Bc447696F431E425D02dd46a8c
data : {1 c758d0418cfd96ed0000 3651d63fc041c58389f4cf0fb3fda66de9a32a0cd2e46abfdfa879c4c58b9834 7ce293e59007112eda7059ed925f5a539ef50eb0997864f24f16007e9f746470}
```
---------------------------
## 2 - FreezeDB

## 2.1 - FreezeBlock
Search in FreezeDB the bloc

```sh
  go run main.go freezeBlock <FreezeDB path> <block number>
```

Returns raw informations store in freezeDB about this bloc.
* header
* hash
* body
* receipts
* diff

Example :
```sh
go run main.go freezeBlock ./.ethereum/geth/chaindata/ancient/ 500
Freeze Data block : 500

header : {2f9dc5dff99590d5f8f742f90e1224eaf0c9c03ba741a0f25f30b5f41abf3e26 1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347 0000000000000000000000000000000000000000 5d6cded585e73c4e322c30c2f782a336316f17dd85a4863b9d838d2d4b8b3008 56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 2 1f4 7a1200 0 5c532d3a 506172697479205465636820417574686f726974790000000000000000000000438098b8726ca83901e4dee8b921fb6f59c410d377377528a57be8ce0f7e63e5550303cb755fd6e836b98576fce895c77549089956360d374ec1f35799e1ffeb01 0000000000000000000000000000000000000000000000000000000000000000 0000000000000000 <nil>}

hash : 0b6e0e5b8c5c9e927af8d56a9e4aa6a7d3170af5979c3c5cb2c65b17dc3c4309

body : {[] []}

receipts : []

diff : 3e9
```
---------------------------
## 3 - Geth tools

## 3.1 - inspect
Same as geth inspect

```sh
  go run main.go inspect <Chaindata path>
```

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<!-- CONTACT -->
## Contact

Thomas Martignon - thomas.martignon@utt.fr