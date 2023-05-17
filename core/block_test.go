package core

import 
(
"github.com/muaj07/transport/types"
"testing"
"time"
"github.com/stretchr/testify/assert"
"bytes"
"fmt"
)

func TestHeader_Encode_Decode(t *testing.T) {
	h:= &Header{
		Version: 1,
		PrevBlock: types.RandomHash(),
		TimeStamp: time.Now().UnixNano(),
		Height: 10,
		Nonce: 110,
	}
	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))
	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)	
	
}

func TestBlock_Encode_Decode(t *testing.T) {
	b:= &Block {
		Header: Header{
			Version: 1,
			PrevBlock: types.RandomHash(),
			TimeStamp: time.Now().UnixNano(),
			Height: 10,
			Nonce: 110,
		},
		Transactions: nil,
	}
	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))
	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
	fmt.Printf("%+v", bDecode)
}

func TestBlockHash(t *testing.T){
	b := &Block{
		Header: Header{
			Version: 1,
			PrevBlock: types.RandomHash(),
			TimeStamp: time.Now().UnixNano(),
			Height: 10,
			Nonce: 110,
		},
		Transactions: nil,
		}

	h := b.Hash()
	fmt.Println(h)
	assert.False(t, h.IsZero())
	}