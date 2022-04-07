package tools

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func FreezerBlockData(freezePath string, blockNb uint64) {
	fmt.Printf("Freeze Data form block : %v\n\n", blockNb)
	
	header := getBlockHeader(freezePath, blockNb)
	fmt.Printf("header : %x\n\n", header)

	hash := getBlockHash(freezePath, blockNb)
	fmt.Printf("hash : %x\n\n", hash)
	
	body := getBlockBody(freezePath, blockNb)
	fmt.Printf("body : %v\n\n", body)

	receipts := getBlockReceipt(freezePath, blockNb)
	fmt.Printf("receipts : %v\n\n", receipts)

	diff := getBlockDiff(freezePath, blockNb)
	fmt.Printf("diff : %x\n\n", diff)
}

// ====================================================================================================

func getBlockHash(freezePath string, blockNumber uint64) []byte {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, freezerHashTable, FreezerNoSnappy[freezerHashTable], true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)

	return data
}

func getBlockDiff(freezePath string, blockNumber uint64) *big.Int {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, freezerDifficultyTable, FreezerNoSnappy[freezerDifficultyTable], true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)

	td := new(big.Int)
	rlp.Decode(bytes.NewReader(data), td)
	return td
}


func getBlockReceipt(freezePath string, blockNumber uint64) types.Receipts {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, freezerReceiptTable, FreezerNoSnappy[freezerReceiptTable], true)
	if err != nil {
		panic(err)
	}
	
	data, _ := freezeDB.Retrieve(blockNumber)

	// Convert the receipts from their storage form to their internal representation
	storageReceipts := []*types.ReceiptForStorage{}
	rlp.DecodeBytes(data, &storageReceipts)
	receipts := make(types.Receipts, len(storageReceipts))
	for i, storageReceipt := range storageReceipts {
		receipts[i] = (*types.Receipt)(storageReceipt)
	}
	return receipts
}

func getBlockBody(freezePath string, blockNumber uint64) types.Body{
	freezeDB, err := rawdb.NewFreezerTable(freezePath, freezerBodiesTable, FreezerNoSnappy[freezerBodiesTable], true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)
	var body types.Body
	rlp.DecodeBytes(data, &body)
	return body
}

func getBlockHeader(freezePath string, blockNumber uint64) types.Header {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, freezerHeaderTable, FreezerNoSnappy[freezerHeaderTable], true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)
	var header types.Header
	rlp.DecodeBytes(data, &header)
	return header
}