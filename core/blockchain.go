package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"log"
)

type Blockchain struct {
	blocks  []*Block
	Current *Block
}

func NewBlockchain(to, message string) *Blockchain {
	coinbase := NewCoinbaseTx(to, message)
	block := NewGenesisBlock(coinbase)
	return &Blockchain{[]*Block{block}, block}
}

func (bc *Blockchain) GetBlocks() []*Block {
	return bc.blocks
}

// var DebugMode = false

func (bc *Blockchain) FindUTXO() map[string]map[int]*TxOut {
	utxo := make(map[string]map[int]*TxOut)
	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			for index, vout := range tx.Vout {
				id := tx.ID.String()
				if utxo[id] == nil {
					utxo[id] = make(map[int]*TxOut)
				}
				utxo[id][index] = vout
			}
			for _, vin := range tx.Vin {
				// FUCK this id, mistake for tx.ID before
				id := vin.TxID.String()
				delete(utxo[id], vin.Vout)
			}
		}
	}
	return utxo
}

func (bc *Blockchain) FindUTXOForKey(pubKeyHash []byte) []*TxOut {
	var utxo []*TxOut
	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			for _, vout := range tx.Vout {
				if vout.IsLockedByKey(pubKeyHash) {
					utxo = append(utxo, vout)
				}
			}
		}
	}
	return utxo
}

// should use map in return, because txout has no index info. mistake to use slice before
func (bc *Blockchain) FindSpendableUTXO(pubKeyHash []byte, amount int) (int, map[string]map[int]*TxOut) {
	utxo := make(map[string]map[int]*TxOut)
	accumulated := 0
	for txid, value := range bc.FindUTXO() {
		for index, txout := range value {
			if txout.IsLockedByKey(pubKeyHash) && accumulated < amount {
				accumulated += txout.Value
				if utxo[txid] == nil {
					utxo[txid] = make(map[int]*TxOut)
				}
				utxo[txid][index] = txout
			}
		}
	}
	return accumulated, utxo
}

func (bc *Blockchain) AddBlock(b *Block) {
	if bc == nil {
		panic("need genesis block")
	}
	b.PrevBlock = bc.blocks[len(bc.blocks)-1].Hash()
	bc.blocks = append(bc.blocks, b)
	bc.Current = b
	return
}

func (bc *Blockchain) FindTransaction(id Hash) (*Transaction, error) {
	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID[:], id[:]) == 0 {
				return tx, nil
			}
		}
	}
	return nil, nil
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey *ecdsa.PrivateKey) {
	prevTxs := make(map[string]*Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[prevTx.ID.String()] = prevTx
	}

	tx.Sign(privKey, prevTxs)
}

func (bc *Blockchain) String() string {
	bt, _ := json.Marshal(bc.blocks)
	return string(bt)
}
