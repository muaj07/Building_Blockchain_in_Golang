package core

import 
(
"github.com/muaj07/transport/types"
"testing"
"time"
"github.com/stretchr/testify/assert"
"github.com/muaj07/transport/crypto"
//"bytes"
//"fmt"
)

func TestSignBlock(t *testing.T) {
	privKey :=crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})
	assert.Nil(t,b.Sign(privKey))
	assert.NotNil(t, b.Signature)
	
}

func TestVerifyBlock(t *testing.T) {
	privKey :=crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t,b.Sign(privKey))
	assert.Nil(t,b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t,b.Verify())
	b.Height=100
	assert.NotNil(t, b.Verify())
}

//Helper function to create a new Block for testing methods
func randomBlock(height uint32, PrevBlockHash types.Hash) *Block{
	header := &Header{
		Version: 1,
		PrevBlockHash: PrevBlockHash,
		TimeStamp: time.Now().UnixNano(),
		Height: height,
		Nonce: 110,
	}
	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature (t *testing.T, height uint32, PrevBlockHash types.Hash) *Block{
privKey := crypto.GeneratePrivateKey()
b:= randomBlock(height, PrevBlockHash)
tx := randomTxWithSignature(t)
b.AddTransaction(tx)
assert.Nil(t, b.Sign(privKey))
return b
}