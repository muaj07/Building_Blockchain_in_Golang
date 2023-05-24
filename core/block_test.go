package core

import 
(
"github.com/muaj07/transport/types"
"testing"
"time"
"github.com/stretchr/testify/assert"
"github.com/muaj07/transport/crypto"
"bytes"
//"fmt"
)

func TestSignBlock(t *testing.T) {
	privKey :=crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})
	assert.Nil(t,b.Sign(privKey))
	assert.NotNil(t, b.Signature)
	
}

func TestVerifyBlock(t *testing.T) {
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t,b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t,b.Verify())

	b.Height=100
	assert.NotNil(t, b.Verify())
}


func TestEncodeDecodeBlock(t *testing.T) {
	b := randomBlock(t, 1, types.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t,b.Encode(NewGobBlockEncoder(buf)))
	bDecode := new(Block)
	assert.Nil(t,bDecode.Decode(NewGobBlockDecoder(buf)))
	assert.Equal(t, b, bDecode)
}


//Helper function to create a new Block for testing methods
func randomBlock(t *testing.T, height uint32, PrevBlockHash types.Hash) *Block{
	privKey := crypto.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version: 1,
		PrevBlockHash: PrevBlockHash,
		TimeStamp: time.Now().UnixNano(),
		Height: height,
		Nonce: 110,
	}
	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t,err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t,err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))
	return b
}
