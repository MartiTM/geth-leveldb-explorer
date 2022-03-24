package tools

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/schollz/progressbar/v3"
)

func CountStorageTree(ldbPath string) {
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}
	bar := progressbar.Default(-1, "Block crowled")
	fmt.Printf("\n")
	total := 0
	headerHash, _ := ldb.Get(headHeaderKey)
	for headerHash != nil {
		var blockHeader types.Header
		blockNb, _ := ldb.Get(append(headerNumberPrefix, headerHash...))
		if blockNb == nil {
			break
		}
		blockHeaderRaw, _ := ldb.Get(append(headerPrefix[:], append(blockNb, headerHash...)...))
		rlp.DecodeBytes(blockHeaderRaw, &blockHeader)

		stateRootNode, _ := ldb.Get(blockHeader.Root.Bytes())
		bar.Add(1)
		if len(stateRootNode) > 0 {
			total++
		}
		headerHash = blockHeader.ParentHash.Bytes()
	}
	fmt.Printf("\nTotal number of state tree : %v\n", total)
} 

// Displays the size of the storage trees of the most recent state tree present in levelDB
func GetStorageTreeSize(ldbPath string) {
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}

	stateRootNode, _ := GetLatestStateTree(ldb)
	
	storageRootNodes := make(chan common.Hash)
	size := make(chan int)
	total := 0
 	
	go GetStorageRootNodes(ldb, stateRootNode, storageRootNodes)

	go func() {
		for storageRoot := range storageRootNodes {
			GetTreeSize(ldb, storageRoot, size)
		}
		close(size)
	}()

	for s := range size {
		total += s
	}

	fmt.Printf("Size in bytes :%v\n", total)
}

// Returns the hash of the most recent state tree 
func GetLatestStateTree(ldb ethdb.Database) (common.Hash, error) {
	headerHash, _ := ldb.Get(headHeaderKey)
	for headerHash != nil {
		var blockHeader types.Header
		blockNb, _ := ldb.Get(append(headerNumberPrefix, headerHash...))
		blockHeaderRaw, _ := ldb.Get(append(headerPrefix[:], append(blockNb, headerHash...)...))
		rlp.DecodeBytes(blockHeaderRaw, &blockHeader)

		stateRootNode, _ := ldb.Get(blockHeader.Root.Bytes())

		if len(stateRootNode) > 0 {
			fmt.Printf("Block number : %x\n", blockHeader.Number)
			fmt.Printf("State root : %x\n", blockHeader.Root)
			return blockHeader.Root, nil
		}
		headerHash = blockHeader.ParentHash.Bytes()
	}
	return common.Hash{}, fmt.Errorf("State tree not found")
}

// Go through the state tree to put in the channel the hashes of the smartcontracts root nodes
func GetStorageRootNodes(ldb ethdb.Database, stateRootNode common.Hash, c chan common.Hash) (error) {
	defer close(c)

	barAcc := progressbar.Default(-1, "Account found")
	fmt.Printf("\n")

	trieDB := trie.NewDatabase(ldb)
	tree, _ := trie.New(stateRootNode, trieDB)

	it := trie.NewIterator(tree.NodeIterator(nil))
	nbAccount := 0
	nbSmartcontract := 0
	for it.Next() {
		var acc snapshot.Account
		barAcc.Add(1)
		nbAccount++

		if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
			panic(err)
		}

		if bytes.Compare(acc.Root, emptyStorageRoot) != 0 {
			nbSmartcontract++
			c <- common.BytesToHash(acc.Root)
		}

	}

	fmt.Printf("\nFinal account number :%v\n", nbAccount)
	fmt.Printf("Final smartcontract number :%v\n", nbSmartcontract)
	
	return nil
}

// Returns in the channel each node size of the tree
func GetTreeSize(ldb ethdb.Database, rootNode common.Hash, s chan int) {
	value, err := ldb.Get(rootNode[:])
	if err != nil {
		return
	}
	
	s <- len(rootNode) + len(value)
	
	var nodes [][]byte
	rlp.DecodeBytes(value, &nodes)

	for _, keyNode := range nodes {
		if len(keyNode) == 0 {
			continue
		}
		GetTreeSize(ldb, common.BytesToHash(keyNode), s)
	}
}