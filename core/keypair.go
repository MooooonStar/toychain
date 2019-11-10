package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	AddressVersion = byte(0x00)
	ChecksumLen    = 4
)

type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewKeyPair() *KeyPair {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return &KeyPair{private, pub}
}

func (k *KeyPair) Address() string {
	pub := sha256.Sum256(k.PublicKey)
	h := ripemd160.New()
	h.Write(pub[:])
	hash := h.Sum(nil)
	payload := append([]byte{AddressVersion}, hash...)
	full := append(payload, checksum(payload)...)
	return base58.Encode(full)
}

func HashPubKey(pub []byte) []byte {
	bts := sha256.Sum256(pub)
	h := ripemd160.New()
	h.Write(bts[:])
	return h.Sum(nil)
}

func checksum(payload []byte) []byte {
	first := sha256.Sum256(payload)
	second := sha256.Sum256(first[:])
	return second[:ChecksumLen]
}

func ValidateAddress(address string) bool {
	pubKeyHash, err := base58.Decode(address)
	if err != nil {
		return false
	}
	actualChecksum := pubKeyHash[len(pubKeyHash)-ChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-ChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
