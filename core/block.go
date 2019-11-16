package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

const (
	// set as a contant in target
	TargetBits = 8
)

type BlockHeader struct {
	Version    uint32
	PrevBlock  Hash
	MerkleRoot Hash
	Timestamp  int64
	Nonce      uint32
	Target     uint32
}

type Block struct {
	*BlockHeader
	Transactions []*Transaction
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes()
}

func (b *Block) Hash() Hash {
	data := bytes.Join(
		[][]byte{
			b.PrevBlock.Bytes(),
			b.HashTransactions().Bytes(),
			IntToHex(b.Timestamp),
			IntToHex(int64(b.Target)),
			IntToHex(int64(b.Nonce)),
		},
		[]byte{},
	)
	return sha256.Sum256(data)
}

func (b *Block) HashTransactions() Hash {
	var txs [][]byte
	for _, tx := range b.Transactions {
		bt, _ := tx.Serialize()
		txs = append(txs, bt)
	}
	tree, err := NewMerkleTree(txs)
	if err != nil {
		log.Panic(err)
	}
	return tree.Root.Hash
}

func NewBlock(transactions []*Transaction, preBlock Hash) *Block {
	header := &BlockHeader{
		PrevBlock: preBlock,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
	}
	block := &Block{header, transactions}
	block.Nonce = ProofOfWork(*block)
	block.MerkleRoot = block.HashTransactions()
	block.Target = TargetBits
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, Hash{})
}

func (b *Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(b)
	return buf.Bytes(), err
}

func (b *Block) Deserialize(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}
