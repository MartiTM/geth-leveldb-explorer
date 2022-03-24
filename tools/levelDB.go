package tools

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/schollz/progressbar/v3"
)

func CountingStorageTrees(ldbPath string) {
	ldb := getLDB(ldbPath)

	stateTrees := getStateTrees(ldb)

	fmt.Printf("\nTotal number of tree state : %v\n", len(stateTrees))
} 

// Displays the size of the most recent storage tree present in levelDB
func LatestStateTreeSize(ldbPath string) {
	ldb := getLDB(ldbPath)

	stateTrees := getStateTrees(ldb)

	fmt.Printf("\nTotal number of tree state : %v\n", len(stateTrees))

	latestStateTree := stateTrees[0]
	fmt.Printf("Latest state tree : %x\n", latestStateTree.blockNumber)
	fmt.Printf("Block number : %x\n", latestStateTree.blockNumber)
	fmt.Printf("State root : %x\n", latestStateTree.stateRoot)
	
	totalSize := getStorageTreeSize(ldb, latestStateTree.stateRoot)

	fmt.Printf("Latest storage trees size : %v bytes\n", totalSize)
}

type stateFound struct {
	blockNumber *big.Int;
	stateRoot common.Hash;
}

func getStateTrees(ldb ethdb.Database) ([]stateFound) {
	var res []stateFound
	bar := progressbar.Default(-1, "Block crowled")
	fmt.Printf("\n")

	headerHash, _ := ldb.Get(headHeaderKey)
	for headerHash != nil {
		var blockHeader types.Header
		blockNb, _ := ldb.Get(append(headerNumberPrefix, headerHash...))
		blockHeaderRaw, _ := ldb.Get(append(headerPrefix[:], append(blockNb, headerHash...)...))
		rlp.DecodeBytes(blockHeaderRaw, &blockHeader)

		stateRootNode, _ := ldb.Get(blockHeader.Root.Bytes())
		
		bar.Add(1)
		if len(stateRootNode) > 0 {
			res = append(res, stateFound{blockHeader.Number, blockHeader.Root})
		}

		headerHash = blockHeader.ParentHash.Bytes()
	}
	return res
}

func getStorageTreeSize(ldb ethdb.Database, stateRootNode common.Hash) int {
	chan_storageRootNodes := make(chan common.Hash)
	
	go getStorageRootNodes(ldb, stateRootNode, chan_storageRootNodes)
	
	chan_nodeSize := make(chan int)

	go func() {
		for storageRoot := range chan_storageRootNodes {
			getTreeSize(ldb, storageRoot, chan_nodeSize)
		}
		close(chan_nodeSize)
	}()

	total := 0

	for s := range chan_nodeSize {
		total += s
	}

	return total
}

// Go through the state tree to put in the channel the hashes of the smartcontracts root nodes
func getStorageRootNodes(ldb ethdb.Database, stateRootNode common.Hash, c chan common.Hash) {
	defer close(c)

	barAcc := progressbar.Default(-1, "Account found")
	fmt.Printf("\n")

	trieDB := trie.NewDatabase(ldb)
	treeState, _ := trie.New(stateRootNode, trieDB)

	it := trie.NewIterator(treeState.NodeIterator(nil))
	nbAccount := 0
	nbSmartcontract := 0
	for it.Next() {
		var acc snapshot.Account
		
		if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
			panic(err)
		}

		barAcc.Add(1)
		nbAccount++
		
		if bytes.Compare(acc.Root, emptyStorageRoot) != 0 {
			nbSmartcontract++
			c <- common.BytesToHash(acc.Root)
		}

	}

	fmt.Printf("\nFinal account number :%v\n", nbAccount)
	fmt.Printf("Final smartcontract number :%v\n", nbSmartcontract)
}

// Returns in the channel each node size of the tree
func getTreeSize(ldb ethdb.Database, rootNode common.Hash, nodeSize chan int) {
	value, err := ldb.Get(rootNode[:])
	if err != nil {
		return
	}
	
	nodeSize <- len(rootNode) + len(value)
	
	var nodes [][]byte
	rlp.DecodeBytes(value, &nodes)

	for _, keyNode := range nodes {
		if len(keyNode) == 0 {
			continue
		}
		getTreeSize(ldb, common.BytesToHash(keyNode), nodeSize)
	}
}

func getLDB(ldbPath string) ethdb.Database {
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}
	fmt.Print("LevelDB ok")
	return ldb
}