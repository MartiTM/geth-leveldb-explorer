package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"

	// "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	// "encoding/binary"
	// "github.com/davecgh/go-spew/spew"
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
	ldbPath := "../.ethereum/geth/chaindata"
	ldb, err := rawdb.NewLevelDBDatabase(ldbPath, 0, 0, "", true)
	if err != nil {
		panic(err)
	}

	stateRootNode, _ := getLastestStateTree(ldb)
	fmt.Printf("State root found :%v\n", stateRootNode) 

	storageRootNodes, _ := getStorageRootNodes(ldb, stateRootNode)

	fmt.Printf("%v\n", len(storageRootNodes))

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
	return common.Hash{}, fmt.Errorf("no state tree found")
}

func getStorageRootNodes(ldb ethdb.Database, stateRootNode common.Hash) ([]common.Hash, error) {
	var storageRootNodes []common.Hash

	trieDB := trie.NewDatabase(ldb)
	state_trie, _ := trie.New(stateRootNode, trieDB)

	it := trie.NewIterator(state_trie.NodeIterator(stateRootNode[:]))
	i:=0
	y:=0
	for it.Next() {
		var acc Account
		i++
		if i%10000 == 0 {
			fmt.Printf("Nombre de compte :%v\n", i)
		}
		if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
			panic(err)
		}
		if bytes.Compare(acc.Root.Bytes(), emptyStorageRoot) != 0 {
			y++
			if y%1000 == 0 {
				fmt.Printf("Smartcontract :%v\n", y)
			}
			storageRootNodes = append(storageRootNodes, acc.Root)
		}
	}

	return storageRootNodes, nil
}
