package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"math/big"
)

/*
	Pay to Public Key Hash (P2PKH),
   | TxIn:scriptSig      | TxOut: scriptPubKey |
	<signature> <pubKey>   OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG
*/
type Transaction struct {
	ID   Hash
	Vin  []*TxIn
	Vout []*TxOut
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1
}

// trimmedCopy copy part of a transaction for signature
func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []*TxIn
	for _, vin := range tx.Vin {
		inputs = append(inputs, &TxIn{vin.TxID, vin.Vout, nil, nil})
	}
	var outputs []*TxOut
	for _, vout := range tx.Vout {
		outputs = append(outputs, &TxOut{vout.Value, vout.PubKeyHash, false})
	}
	return &Transaction{tx.ID, inputs, outputs}
}

func (tx *Transaction) Sign(priv *ecdsa.PrivateKey, preTxs map[string]*Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}
	for _, vin := range tx.Vin {
		if _, exist := preTxs[vin.TxID.String()]; !exist {
			return errors.New("wrong input")
		}
	}

	txCopy := tx.trimmedCopy()
	for i, txin := range txCopy.Vin {
		preTx := preTxs[txin.TxID.String()]
		txCopy.Vin[i].Signature = nil
		txCopy.Vin[i].PubKey = preTx.Vout[txin.Vout].PubKeyHash

		txBytes, _ := txCopy.Serialize()
		payload := fmt.Sprintf("%x\n", txBytes)
		r, s, err := ecdsa.Sign(rand.Reader, priv, []byte(payload))
		if err != nil {
			return err
		}
		tx.Vin[i].Signature = append(r.Bytes(), s.Bytes()...)
	}
	return nil
}

func (tx *Transaction) Hash() Hash {
	txCopy := *tx
	txCopy.ID = Hash{}
	bt, _ := txCopy.Serialize()
	return sha256.Sum256(bt)
}

func (tx *Transaction) Verify(preTxs map[string]*Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}
	for _, txin := range tx.Vin {
		if preTxs[txin.TxID.String()] == nil {
			return errors.New("wrong input")
		}
	}
	txCopy := tx.trimmedCopy()
	for i, txin := range tx.Vin {
		prevTx := preTxs[txin.TxID.String()]
		txCopy.Vin[i].Signature = nil
		txCopy.Vin[i].PubKey = prevTx.Vout[txin.Vout].PubKeyHash

		var r, s *big.Int
		sigLen := len(txin.Signature)
		r.SetBytes(txin.Signature[:(sigLen / 2)])
		s.SetBytes(txin.Signature[(sigLen / 2):])

		var x, y *big.Int
		keyLen := len(txin.PubKey)
		x.SetBytes(txin.PubKey[:(keyLen / 2)])
		y.SetBytes(txin.PubKey[(keyLen / 2):])

		txBytes, _ := txCopy.Serialize()
		payload := fmt.Sprintf("%x\n", txBytes)
		pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

		if !ecdsa.Verify(pub, []byte(payload), r, s) {
			return errors.New("invalid transation")
		}
	}
	return nil
}

const (
	// 挖矿奖励恒定为100
	Prize = 100
)

func NewCoinbaseTx(to, data string) *Transaction {
	txin := &TxIn{Hash{}, -1, nil, []byte(data)}
	txout := NewTxOut(Prize, to)
	tx := &Transaction{Hash{}, []*TxIn{txin}, []*TxOut{txout}}
	tx.ID = tx.Hash()
	return tx
}

func NewTransaction(to string, amount int, key *KeyPair, bc *Blockchain) (*Transaction, error) {
	pubKeyHash := HashPubKey(key.PublicKey)
	acc, validOuts := bc.FindSpendableUTXO(pubKeyHash, amount)
	if acc < amount {
		return nil, fmt.Errorf("insufficinent balance: %v", acc)
	}

	var inputs []*TxIn
	for txid, txOuts := range validOuts {
		id, _ := HashFromString(txid)
		for index, _ := range txOuts {
			input := &TxIn{id, index, nil, key.PublicKey}
			inputs = append(inputs, input)
		}
	}

	var outputs []*TxOut
	from := fmt.Sprintf("%s", key.Address())
	outputs = append(outputs, NewTxOut(amount, to))
	if acc > amount {
		outputs = append(outputs, NewTxOut(acc-amount, from))
	}

	tx := &Transaction{Hash{}, inputs, outputs}
	tx.ID = tx.Hash()

	bc.SignTransaction(tx, key.PrivateKey)
	return tx, nil
}

func (tx *Transaction) Serialize() ([]byte, error) {
	if tx == nil {
		return nil, nil
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(tx)
	return buf.Bytes(), err
}

func (tx *Transaction) Deserialize(data []byte) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(tx)
}
