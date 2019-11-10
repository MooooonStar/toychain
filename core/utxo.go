package core

import (
	"bytes"

	"github.com/mr-tron/base58"
)

type TxIn struct {
	TxID      Hash // 输入指向的交易ID
	Vout      int  // 所指向交易创建的TxOut索引
	Signature []byte
	PubKey    []byte
}

func (in *TxIn) SameKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

type TxOut struct {
	Value      int
	PubKeyHash []byte
	spent      bool
}

// Lock means put a address to TxOut
func (out *TxOut) Lock(address string) error {
	pub, err := base58.Decode(address)
	if err != nil {
		return err
	}
	out.PubKeyHash = pub[1 : len(pub)-4]
	return nil
}

func (out *TxOut) IsLockedByKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTxOut(value int, address string) *TxOut {
	txo := &TxOut{value, nil, false}
	txo.Lock(address)
	return txo
}
