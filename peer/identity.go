package peer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/bits"

	"github.com/perlin-network/noise/crypto/blake2b"
)

// ID is an identity of nodes, using its public key hash and network address.
type ID struct {
	// Address is the network address of the peer
	Address string
	// ID is the computed hash of the public key
	ID []byte
	// PublicKey of the peer (we no longer use the public key as the peer ID, but use it to verify messages)
	PublicKey []byte
}

// CreateID is a factory function creating ID.
func CreateID(address string, publicKey []byte) ID {
	return ID{Address: address, PublicKey: publicKey, ID: blake2b.New().HashBytes(publicKey)}
}

// String returns the identity address and public key.
func (id ID) String() string {
	return fmt.Sprintf("ID{Address: %v, ID: %v}", id.Address, id.ID)
}

// Equals determines if two peer IDs are equal to each other based on the contents of their public keys.
func (id ID) Equals(other ID) bool {
	return bytes.Equal(id.ID, other.ID)
}

// Less determines if this peer ID's public key is less than other ID's public key.
func (id ID) Less(other interface{}) bool {
	if other, is := other.(ID); is {
		return bytes.Compare(id.ID, other.ID) == -1
	}
	return false
}

// PublicKeyHex generates a hex-encoded string of public key hash of this given peer ID.
func (id ID) PublicKeyHex() string {
	return hex.EncodeToString(id.PublicKey)
}

// Xor performs XOR (^) over another peer ID's public key.
func (id ID) Xor(other ID) ID {
	result := make([]byte, len(id.PublicKey))

	for i := 0; i < len(id.PublicKey) && i < len(other.PublicKey); i++ {
		result[i] = id.PublicKey[i] ^ other.PublicKey[i]
	}
	return ID{Address: id.Address, PublicKey: result}
}

// XorID performs XOR (^) over another peer ID's public key hash.
func (id ID) XorID(other ID) ID {
	result := make([]byte, len(id.ID))

	for i := 0; i < len(id.ID) && i < len(other.ID); i++ {
		result[i] = id.ID[i] ^ other.ID[i]
	}
	return ID{Address: id.Address, ID: result}
}

// PrefixLen returns the number of prefixed zeros in a peer ID.
func (id ID) PrefixLen() int {
	for i, b := range id.ID {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(id.ID)*8 - 1
}
