package main

import (
	"encoding/json"
	"log"

	"github.com/MooooonStar/toychain/core"
)

func main() {
	addr := "17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv"

	key := core.NewKeyPair()
	to, message := key.Address(), "Long Live The Bitcoin"
	bc := core.NewBlockchain(to, message)

	key1 := core.NewKeyPair()
	tx1, err := core.NewTransaction(key1.Address(), 6, key, bc)
	if err != nil {
		log.Panic(err)
	}

	block1 := core.NewBlock([]*core.Transaction{tx1}, bc.Current.Hash())
	nonce1 := core.ProofOfWork(*block1)
	block1.Nonce = nonce1
	bc.AddBlock(block1)

	tx2, err := core.NewTransaction(addr, 2, key1, bc)
	if err != nil {
		log.Panic(err)
	}
	tx3, err := core.NewTransaction(addr, 3, key, bc)
	if err != nil {
		log.Panic(err)
	}
	block2 := core.NewBlock([]*core.Transaction{tx2, tx3}, bc.Current.Hash())
	nonce2 := core.ProofOfWork(*block2)
	block2.Nonce = nonce2
	bc.AddBlock(block2)

	bt, _ := json.Marshal(bc.GetBlocks())
	log.Println(string(bt))
}
