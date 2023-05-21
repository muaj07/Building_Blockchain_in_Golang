package core

import(
	"io"
	"encoding/gob"
	"crypto/elliptic"
)


type Encoder [T any] interface{
	Encode(T) error

}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder (w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{
		w: w}
}

func (g *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(g.w).Encode(tx)
	
}