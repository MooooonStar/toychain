package main

import (
	"log"

	"github.com/MooooonStar/toychain/core"
	"github.com/hokaccha/go-prettyjson"
)

func main() {
	// 创世块 T
	addr := "17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv"
	message := "Long Live The Bitcoin"
	bc := core.NewBlockchain(addr, message)

	// A挖到了第一个块
	keyA := core.NewKeyPair()
	tx0 := core.NewCoinbaseTx(keyA.Address(), "A got it")
	block1 := core.NewBlock([]*core.Transaction{tx0}, bc.Current.Hash())
	nonce1 := core.ProofOfWork(*block1)
	block1.Nonce = nonce1
	bc.AddBlock(block1)

	// B挖到了第二个块，并且A向C矿工转账了6
	keyB, keyC := core.NewKeyPair(), core.NewKeyPair()
	tx1 := core.NewCoinbaseTx(keyB.Address(), "B got it")
	tx2, err := core.NewTransaction(keyC.Address(), 6, keyA, bc)
	if err != nil {
		log.Fatal(err)
	}
	block2 := core.NewBlock([]*core.Transaction{tx1, tx2}, bc.Current.Hash())
	nonce2 := core.ProofOfWork(*block2)
	block1.Nonce = nonce2
	bc.AddBlock(block2)

	//D挖到了第三个块, C->D 2, A->T 4
	keyD := core.NewKeyPair()
	tx3 := core.NewCoinbaseTx(keyD.Address(), "D got it")
	tx4, err := core.NewTransaction(keyD.Address(), 2, keyC, bc)
	if err != nil {
		log.Fatal(err)
	}
	tx5, err := core.NewTransaction(keyB.Address(), 4, keyA, bc)
	if err != nil {
		log.Fatal(err)
	}
	block3 := core.NewBlock([]*core.Transaction{tx3, tx4, tx5}, bc.Current.Hash())
	nonce3 := core.ProofOfWork(*block3)
	block3.Nonce = nonce3
	bc.AddBlock(block3)

	bt, _ := prettyjson.Marshal(bc.GetBlocks())
	log.Println(string(bt))

	for _, key := range []*core.KeyPair{keyA, keyB, keyC, keyD} {
		_, err := core.NewTransaction(addr, 1000, key, bc)
		if err != nil {
			log.Println("show balance:", err)
		}
	}
}
