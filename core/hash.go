package core

import (
	"encoding/hex"
	"errors"
)

// Hash is 256-bit
type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func HashFromString(hash string) (h Hash, err error) {
	bt, err := hex.DecodeString(hash)
	if err != nil {
		return
	}
	if len(bt) != 32 {
		return Hash{}, errors.New("invalid hash length")
	}
	copy(h[:], bt)
	return
}
