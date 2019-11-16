package core

import (
	"crypto/sha256"
	"math"
	"math/big"
)

func ProofOfWork(b Block) uint32 {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Target))

	bb := b
	bb.Nonce = uint32(0)
	for bb.Nonce < math.MaxUint32 {
		hashInt := new(big.Int)
		hash := bb.Hash()
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			break
		}
		bb.Nonce++
	}

	return bb.Nonce
}

func CheckProofOfWork(b Block) bool {
	var hashInt big.Int
	h := b.Hash()
	hash := sha256.Sum256(h[:])
	hashInt.SetBytes(hash[:])

	target := big.NewInt(1)
	target.Lsh(target, uint(256-b.Target))
	return hashInt.Cmp(target) == -1
}
