package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/schollz/progressbar/v3"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	// databaseVerisionKey tracks the current database version.
	databaseVerisionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest know header's hash.
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest know full block's hash.
	headBlockKey = []byte("LastBlock")

	// headFastBlockKey tracks the latest known incomplete block's hash during fast sync.
	headFastBlockKey = []byte("LastFast")

	// fastTrieProgressKey tracks the number of trie entries imported during fast sync.
	fastTrieProgressKey = []byte("TrieSync")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
	headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix  = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits

	preimagePrefix = []byte("secure-key-")      // preimagePrefix + hash -> preimage
	configPrefix   = []byte("ethereum-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	emptyStorageRoot, _ = hex.DecodeString("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash
	CodeHash []byte
}

func main() {
	// ldbPath := "../.ethereum/geth/chaindata"
	ldbPath := "../.ethereum-testnet/goerli/geth/chaindata"
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}

	stateRootNode, _ := getLastestStateTree(ldb)
	fmt.Printf("State root found :%v\n", stateRootNode)
	
	storageRootNodes := make(chan common.Hash)
	size := make(chan int)
	defer close(size)
 	
	go getStorageRootNodes(ldb, stateRootNode, storageRootNodes)

	go countSize(size)

	for storageRoot := range storageRootNodes {
		go getTreeSize(ldb, storageRoot, size)
	}
}

func countSize(size chan int) {
	sum := 0

	for s := range size {
		sum += s
	}
	fmt.Printf("size in byte :%v\n", sum)
}

func getLastestStateTree(ldb ethdb.Database) (common.Hash, error) {
	headerHash, _ := ldb.Get(headHeaderKey)
	for headerHash != nil {
		var b types.Header
		blockNb, _ := ldb.Get(append(headerNumberPrefix, headerHash...))
		blockHeader, _ := ldb.Get(append(headerPrefix[:], append(blockNb, headerHash...)...))
		rlp.DecodeBytes(blockHeader, &b)

		stateRootNode, _ := ldb.Get(b.Root.Bytes())

		if len(stateRootNode) > 0 {
			return b.Root, nil
		}
		headerHash = b.ParentHash.Bytes()
	}
	return common.Hash{}, fmt.Errorf("State tree not found")
}

func getStorageRootNodes(ldb ethdb.Database, stateRootNode common.Hash, c chan common.Hash) (error) {
	defer close(c)

	barAcc := progressbar.Default(-1, "Account found")
	fmt.Printf("\n")

	trieDB := trie.NewDatabase(ldb)
	tree, _ := trie.New(stateRootNode, trieDB)

	it := trie.NewIterator(tree.NodeIterator(stateRootNode[:]))
	nbAccount := 0
	nbSmartcontract := 0
	for it.Next() {
		var acc Account
		barAcc.Add(1)
		nbAccount++

		if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
			panic(err)
		}

		if bytes.Compare(acc.Root.Bytes(), emptyStorageRoot) != 0 {
			nbSmartcontract++
			c <- acc.Root
		}

	}

	fmt.Printf("Final account number :%v\n", nbAccount)
	fmt.Printf("Final smartcontract number :%v\n", nbSmartcontract)
	
	return nil
}

func getTreeSize(ldb ethdb.Database, rootNode common.Hash, s chan int) {
	value, err := ldb.Get(rootNode[:])
	if err != nil {
		return
	}
	
	s <- len(rootNode) + len(value)
	
	var nodes [][]byte
	rlp.DecodeBytes(value, &nodes)

	// end of tree
	if len(nodes) == 2 {
		return
	}

	for _, keyNode := range nodes {
		if len(keyNode) == 0 {
			continue
		}
		getTreeSize(ldb, common.BytesToHash(keyNode), s)
	}
}
