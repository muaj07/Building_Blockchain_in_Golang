package core
import (
	"encoding/gob"
	"crypto/elliptic"
	"io"
)


type Decoder [T any] interface{
	Decode( T) error

}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder (r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{
		r:r}
}

func (g *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(g.r).Decode(tx)
	
}