package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/muaj07/transport/types"
	"math/big"
	
)

type PublicKey struct{
	Key *ecdsa.PublicKey
}

type PrivateKey struct{
	Key *ecdsa.PrivateKey
}

type Signature struct{
	R, S *big.Int
	
}

// Generate and return Private Key
func GeneratePrivateKey () PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err !=nil {
		panic(err)
	}
	return PrivateKey {
		Key: key,
	}
}
// Generate and return Public key from the Private Key
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey {
		Key: &k.Key.PublicKey,
	}
}
// Helper method/function
func (k PublicKey) ToSlice () []byte{
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

// Generate the Address from the Public Key
func (k PublicKey) Address () types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.AddressFromBytes(h[len(h)-20:])
}

// Sign the data using the Private key and return the Signature
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
}


// Verify if the data is signed by the right private key
func (sig *Signature) Verify(pubkey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubkey.Key, data, sig.R, sig.S)
}