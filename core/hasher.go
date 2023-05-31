package core

import (
	"github.com/muaj07/transport/types"
	"crypto/sha256"
)

type Hasher [T any] interface {
	Hash(T) types.Hash	
}

type BlockHasher struct {}

// Hash calculates the hash of a header block.
func (BlockHasher) Hash(b *Header) types.Hash {
    // Convert the header block to bytes and hash them with SHA256.
    h := sha256.Sum256(b.Bytes())
    // Return the resulting hash as a types.Hash.
    return types.Hash(h)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}