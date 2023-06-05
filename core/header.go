package core

import (
//"encoding/binary"
"github.com/muaj07/transport/types"
"bytes"
//"crypto/sha256"
//"fmt"
//"github.com/muaj07/transport/crypto"
"encoding/gob"
//"time"
)

type Header struct {
	Version uint32
	DataHash types.Hash
	PrevBlockHash types.Hash
	TimeStamp int64
	Height uint32
	Nonce uint64
}

// Bytes returns the bytes of a Header by encoding it with gob.
func (h *Header) Bytes() []byte {
    // Create a new bytes buffer.
    buf := &bytes.Buffer{}

    // Create a new gob encoder that writes to the buffer.
    enc := gob.NewEncoder(buf)

    // Encode the Header into the buffer.
    enc.Encode(h)

    // Return the bytes of the buffer.
    return buf.Bytes()
}
