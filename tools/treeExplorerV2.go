package tools

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func StateAndStorageTreesV2(ldbPath string) {
	var wg sync.WaitGroup

	p := mpb.New(mpb.WithWaitGroup(&wg))

	ldb := getLDB(ldbPath)

	stateRoot := getLatestStateRootV2(ldb)

	storageRoot := make(chan []byte)

	countSmartcontract := 0

	var (
		stateTree stat
		storageTree stat

		stateLeaf stat
		storageLeaf stat
	) 

	start := time.Now()

	accBar := p.AddBar(180000000,
		mpb.PrependDecorators(
			decor.Name("Account found :	"),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
		),
	)
	
	stateNodeBar := p.AddBar(-1,
			mpb.BarWidth(0),
			mpb.PrependDecorators(
				decor.Name("State nodes : 	"),
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidthR),
			),
	)
	
	storageNodeBar := p.AddBar(-1,
			mpb.BarWidth(0),
			mpb.PrependDecorators(
				decor.Name("Storage nodes :	"),
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidthR),
			),
	)
	
	goRoutineNodeBar := p.AddBar(-1,
			mpb.BarWidth(0),
			mpb.PrependDecorators(
				decor.Name("GoRoutine : 	"),
				decor.CountersNoUnit("%d / %d", decor.WCSyncWidthR),
			),
	)

	wg.Add(1)
	// explore the latest state tree
	go func(){
		defer wg.Done()
		defer close(storageRoot)

		exploreTreeV2(ldb, stateRoot, 
			// check if the node is a leaf. In the leaf extract the account and send the storage root in a channel
			func(node [][]byte, size int) bool {
				stateNodeBar.Increment()
				stateTree.Add(common.StorageSize(size))

				if len(node) == 2 {
					// account check
					var acc snapshot.Account
					if err := rlp.DecodeBytes(node[1], &acc); err == nil {
						accBar.Increment()
						stateLeaf.Add(common.StorageSize(size))

						storageRoot <- acc.Root
						return true
					}
				}
				return false
		})
		p.Abort(accBar, false)
		p.Abort(stateNodeBar, false)
	}()
	
	wg.Add(1)
	// explore each storage tree
	go func(){
		defer wg.Done()
		// handle storage root
		for root := range storageRoot {
			// check if it is a smartcontract
			if bytes.Compare(root, emptyStorageRoot) == 0 {
				continue
			}
			countSmartcontract++
			go func(){
				goRoutineNodeBar.Increment()
				exploreTreeV2(ldb, root,
					func(node [][]byte, size int) bool {
						storageNodeBar.Increment()
						storageTree.Add(common.StorageSize(size))
		
						if len(node) == 2 {
							storageLeaf.Add(common.StorageSize(size))
						}
						return false
				})
				goRoutineNodeBar.IncrBy(-1)
			}()
		}
		p.Abort(storageNodeBar, false)
		p.Abort(goRoutineNodeBar, false)
	}()


	p.Wait()

	fmt.Printf("time : %v \n", time.Now().Sub(start))
	fmt.Printf("state tree :\n")
	fmt.Printf("  - nodes :: number : %v / size : %v\n", stateTree.Count(), stateTree.Size())
	fmt.Printf("  - leafs :: number : %v / size : %v\n", stateLeaf.Count(), stateLeaf.Size())
	fmt.Printf("storage tree :\n")
	fmt.Printf("  - nodes :: number : %v / size : %v\n", storageTree.Count(), storageTree.Size())
	fmt.Printf("  - leafs :: number : %v / size : %v\n", storageLeaf.Count(), storageLeaf.Size())
	fmt.Printf("Account : %v\n", int(stateLeaf.count) - countSmartcontract)
	fmt.Printf("Smartcontract : %v\n", countSmartcontract)
}

func getLatestStateRootV2(ldb ethdb.Database) []byte {
	// latest header hash in LevelDB
	headerHash, _ := ldb.Get(headHeaderKey)
	for headerHash != nil {
		blockNb, err := ldb.Get(append(headerNumberPrefix, headerHash...))
		if blockNb == nil || err != nil {
			panic("No block number")
		}
		blockHeaderRaw, err := ldb.Get(append(headerPrefix[:], append(blockNb, headerHash...)...))
		if blockHeaderRaw == nil || err != nil {
			panic("No block Header")
		}
		var blockHeader types.Header
		rlp.DecodeBytes(blockHeaderRaw, &blockHeader)

		stateRootNode, _ := ldb.Get(blockHeader.Root.Bytes())

		if len(stateRootNode) > 0 {
			fmt.Printf("Latest state tree in the block : %v\n", blockHeader.Number)
			return blockHeader.Root.Bytes()
		}

		headerHash = blockHeader.ParentHash.Bytes()
	}
	panic("No state tree")
}

type isLeaf func(node [][]byte, size int) bool

func exploreTreeV2(ldb ethdb.Database, key []byte, f isLeaf) {
	data, err := ldb.Get(key)
	if err != nil {
		return
	}

	var node [][]byte
	if err := rlp.DecodeBytes(data, &node); err != nil {
		return
	}

	// is leaf ?
	if f(node, len(key)+len(data)) {
		return
	}

	// explore next nodes
	for _, keyNode := range node {
		if len(keyNode) == 0 {
			continue
		}
		exploreTreeV2(ldb, keyNode, f)
	}
}