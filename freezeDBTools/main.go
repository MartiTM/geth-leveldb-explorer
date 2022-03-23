package freezeDBTools

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func GetBlockData(freezePath string, blockNb uint64) {
	header := GetBlockHeader(freezePath, blockNb)
	fmt.Printf("header : %x\n\n", header)

	body := GetBlockBody(freezePath, blockNb)
	fmt.Printf("body : %x\n\n", body)

	receipts := GetBlockReceipt(freezePath, blockNb)
	fmt.Printf("receipts : %x\n\n", receipts)

	diff := GetBlockDiff(freezePath, blockNb)
	fmt.Printf("diff : %x\n\n", diff)

	hash := GetBlockHash(freezePath, blockNb)
	fmt.Printf("hash : %x\n\n", hash)
}

func GetBlockHash(freezePath string, blockNumber uint64) []byte {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, "hashes", true, true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)

	return data
}

func GetBlockDiff(freezePath string, blockNumber uint64) *big.Int {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, "diffs", true, true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)

	td := new(big.Int)
	rlp.Decode(bytes.NewReader(data), td)
	return td
}


func GetBlockReceipt(freezePath string, blockNumber uint64) types.Receipts {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, "receipts", false, true)
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

func GetBlockBody(freezePath string, blockNumber uint64) types.Body{
	freezeDB, err := rawdb.NewFreezerTable(freezePath, "bodies", false, true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)
	var body types.Body
	rlp.DecodeBytes(data, &body)
	return body
}

func GetBlockHeader(freezePath string, blockNumber uint64) types.Header {
	freezeDB, err := rawdb.NewFreezerTable(freezePath, "headers", false, true)
	if err != nil {
		panic(err)
	}

	data, _ := freezeDB.Retrieve(blockNumber)
	var header types.Header
	rlp.DecodeBytes(data, &header)
	return header
}