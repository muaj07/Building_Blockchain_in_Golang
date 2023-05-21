package core

import (
	"testing"
	"github.com/muaj07/transport/crypto"
	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestSignTransaction(t *testing.T) {
	privKey :=crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("First test trx."),
	}
	assert.Nil(t,tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
	
}

func TestVerifyTransaction(t *testing.T) {
	privKey :=crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("First test trx."),
	}
	assert.Nil(t,tx.Sign(privKey))
	assert.Nil(t,tx.Verify())
	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()
	assert.NotNil(t,tx.Verify())
}


func TestTxEncodeDecode(t *testing.T){
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t,tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t,txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t,tx,txDecoded)
}
func randomTxWithSignature (t *testing.T) *Transaction{
	privKey := crypto.GeneratePrivateKey()
	tx:= &Transaction{
		Data: []byte("First test trx with Signature"),
	}
	assert.Nil(t, tx.Sign(privKey))
	return tx
}