package levelDbTree

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/schollz/progressbar/v3"
)

func OldStateExplorer(ldb ethdb.Database, stateRootNode common.Hash) {

	trieDB := trie.NewDatabase(ldb)
	tree, _ := trie.New(stateRootNode, trieDB)

	barAcc := progressbar.Default(-1, "Account found")


	it := trie.NewIterator(tree.NodeIterator(stateRootNode[:]))
	i := 0
	y := 0
	for it.Next() {
		var acc Account

		if err := rlp.DecodeBytes(it.Value, &acc); err != nil {
			panic(err)
		}

		i++
		barAcc.Add(1)
		if bytes.Compare(acc.Root.Bytes(), EmptyStorageRoot) != 0 {
			y++
		}
	}

	fmt.Printf("Nombre de compte Final:%v\n", i)
	fmt.Printf("Smartcontract Final:%v\n", y)
}

func NewStateExplorer(ldb ethdb.Database, stateRootNode common.Hash) {
	x := 0
	z := 0
	barAcc := progressbar.Default(-1, "Account found")


	accounts := make(chan []byte)
	go func() {
		explore(ldb, stateRootNode, accounts)
		close(accounts)
	}()

	for data := range accounts {
		var acc Account
		
		if err := rlp.DecodeBytes(data, &acc); err != nil {
			continue
		}
		x++
		barAcc.Add(1)

		if bytes.Compare(acc.Root.Bytes(), EmptyStorageRoot) != 0 {
			z++
		}

	}
	fmt.Printf("Nombre de compte :%v\n", x)
	fmt.Printf("Smartcontract :%v\n", z)
}

func explore(ldb ethdb.Database, rootNode common.Hash, accounts chan []byte) {
	value, err := ldb.Get(rootNode[:])
	if err != nil {
		return
	}

	var nodes [][]byte
	rlp.DecodeBytes(value, &nodes)
	
	// end of tree
	if len(nodes) == 2 {
		accounts <- nodes[1]
	}

	for _, keyNode := range nodes {
		if len(keyNode) == 0 {
			continue
		}

		explore(ldb, common.BytesToHash(keyNode), accounts)
	}
}
