package core

import(
	"io"
	"encoding/gob"
	"crypto/elliptic"
)


// use GOB encoding for now to bootstrapp the project
// will switch to Protobuf as default encoding/decoding in future

type Encoder [T any] interface{
	Encode(T) error

}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder (w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{
		w: w}
}

func (g *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(g.w).Encode(tx)
	
}

// Encoder for the block

type GobBlockEncoder struct {
	w io.Writer
}

func NewGobBlockEncoder (w io.Writer) *GobBlockEncoder {
	return &GobBlockEncoder{
		w: w,
	}
}

func (enc *GobBlockEncoder) Encode(b *Block) error {
	return gob.NewEncoder(enc.w).Encode(b)
	
}

func init(){
	gob.Register(elliptic.P256())
}