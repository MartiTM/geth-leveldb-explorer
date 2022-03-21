package levelDbTree

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/schollz/progressbar/v3"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
)

func RunTreeSize(ldb ethdb.Database, storageRootNodes chan common.Hash, size chan int) {
	for storageRoot := range storageRootNodes {
		GetTreeSize(ldb, storageRoot, size)
	}
	close(size)
}

func GetLastestStateTree(ldb ethdb.Database) (common.Hash, error) {
	headerHash, _ := ldb.Get(HeadHeaderKey)
	for headerHash != nil {
		var b types.Header
		blockNb, _ := ldb.Get(append(HeaderNumberPrefix, headerHash...))
		blockHeader, _ := ldb.Get(append(HeaderPrefix[:], append(blockNb, headerHash...)...))
		rlp.DecodeBytes(blockHeader, &b)

		stateRootNode, _ := ldb.Get(b.Root.Bytes())

		if len(stateRootNode) > 0 {
			return b.Root, nil
		}
		headerHash = b.ParentHash.Bytes()
	}
	return common.Hash{}, fmt.Errorf("State tree not found")
}

func GetStorageRootNodes(ldb ethdb.Database, stateRootNode common.Hash, c chan common.Hash) (error) {
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

		if bytes.Compare(acc.Root.Bytes(), EmptyStorageRoot) != 0 {
			nbSmartcontract++
			c <- acc.Root
		}

	}

	fmt.Printf("Final account number :%v\n", nbAccount)
	fmt.Printf("Final smartcontract number :%v\n", nbSmartcontract)
	
	return nil
}

func GetTreeSize(ldb ethdb.Database, rootNode common.Hash, s chan int) {
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
		GetTreeSize(ldb, common.BytesToHash(keyNode), s)
	}
}
