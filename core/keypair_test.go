package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert.False(t, ValidateAddress("1mvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"))
	assert.True(t, ValidateAddress("1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"))
	assert.True(t, ValidateAddress("3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy"))

	addr := NewKeyPair().Address()
	assert.True(t, ValidateAddress(addr))

}
