package core

import (
	"encoding/hex"
	"encoding/json"
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

func (h Hash) MarshalJSON() ([]byte, error) {
	str := hex.EncodeToString(h[:])
	return json.Marshal(str)
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	bt, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	copy(h[:], bt[:32])
	return nil
}
