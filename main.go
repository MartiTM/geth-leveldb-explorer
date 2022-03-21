package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/MartiTM/geth-leveldb-explorer/levelDbTree"
)

func main()  {
}

func current()  {
	ldbPath := "../.ethereum/geth/chaindata"
	// ldbPath := "../.ethereum-testnet/goerli/geth/chaindata"
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}

	stateRootNode, _ := levelDbTree.getLastestStateTree(ldb)
	fmt.Printf("State root found :%v\n", stateRootNode)
	
	storageRootNodes := make(chan common.Hash)
	size := make(chan int)
	total := 0
 	
	go levelDbTree.GetStorageRootNodes(ldb, stateRootNode, storageRootNodes)

	go levelDbTree.runTreeSize(ldb, storageRootNodes, size)

	for s := range size {
		total += s
	}

	fmt.Printf("Size in byte :%v\n", total)
}

func test()  {
	// ldbPath := "../.ethereum/geth/chaindata"
	ldbPath := "../.ethereum-testnet/goerli/geth/chaindata"
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}

	stateRootNode, _ := levelDbTree.getLastestStateTree(ldb)
	fmt.Printf("State root found :%v\n", stateRootNode)

	fmt.Printf("sans lib :\n")
	levelDbTree.newStateExplorer(ldb, stateRootNode)

	// fmt.Printf("avec lib :\n")
	// getStorageRootNodesTest(ldb, stateRootNode)
}