package crypto

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)
func TestKeyPair_Sign_Verify_Success(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	//address := pubKey.Address()
	msg := []byte("Hello Enoda!")
	sig, err := privKey.Sign(msg)
	assert.Nil(t,err)
	assert.True(t,sig.Verify(pubKey,msg))
	//fmt.Println(address)
	fmt.Println(sig)
	
}

func TestKeyPair_Sign_Verify_Fail(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Hello Enoda!")
	sig, err := privKey.Sign(msg)
	assert.Nil(t,err)
	otherPrivKey := GeneratePrivateKey()
	otherPubKey := otherPrivKey.PublicKey()
	assert.False(t,sig.Verify(otherPubKey,msg))
	assert.False(t,sig.Verify(pubKey, []byte("Not the original data")))
}