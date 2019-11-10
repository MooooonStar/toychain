package core

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUXTO(t *testing.T) {
	addr := "17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv"

	key := NewKeyPair()
	to, message := key.Address(), "Long Live The Bitcoin"
	bc := NewBlockChain(to, message)

	key1 := NewKeyPair()
	tx1, err := NewTransaction(key, key1.Address(), 6, bc)
	assert.Nil(t, err)
	block1 := NewBlock([]*Transaction{tx1}, bc.Current.Hash())
	nonce1 := ProofOfWork(*block1)
	block1.Nonce = nonce1
	// assert.True(t, CheckProofOfWork(*block1))
	bc.AddBlock(block1)

	tx2, err := NewTransaction(key1, addr, 2, bc)
	assert.Nil(t, err)
	tx3, err := NewTransaction(key, addr, 3, bc)
	assert.Nil(t, err)
	block2 := NewBlock([]*Transaction{tx2, tx3}, bc.Current.Hash())
	nonce2 := ProofOfWork(*block2)
	block2.Nonce = nonce2
	//assert.True(t, CheckProofOfWork(*block2))
	bc.AddBlock(block2)

	bt, err := json.Marshal(bc.blocks)
	assert.Nil(t, err)
	log.Println(string(bt))
}
