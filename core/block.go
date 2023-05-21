package core

import (
//"encoding/binary"
"github.com/muaj07/transport/types"
//"io"
"bytes"
//"crypto/sha256"
"fmt"
"github.com/muaj07/transport/crypto"
"encoding/gob"
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

type Block struct {
	*Header
	Transactions []Transaction
	Validator crypto.PublicKey //Pubkey of the Block validator
	Signature *crypto.Signature //Block signtaure
	// cached version of the Header hash
	hash types.Hash
}

//Constructor function for Block struct
// Only two fields are included in the signature

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header: h,
		Transactions: txx,
	}	
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

// Sign the Block using the "Sign" method from
// the "keypair.go" file in the crypto Folder
// The "Sign method" returns signture (or error), which is assigned to sig
// and then stored in "b.signature", while the Publickey of the Validator
// is stored in "b.Validator"

func (b *Block) Sign(privkey crypto.PrivateKey) error{
	sig, err := privkey.Sign(b.Header.Bytes())
	if err!=nil{
		return err
	}
	b.Signature = sig
	b.Validator = privkey.PublicKey()
	return nil
}


// Verify the Block by checking
// if the "Validator signature" exist and then
// if the transaction is signed by the correct (i.e. of the Validator) Private key

func (b *Block) Verify() error {
	if b.Signature == nil{
		return fmt.Errorf ("Block signature is Nil")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf ("Invalid Block signature")
	}
	// Verify all the transactions in the block
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err!=nil {
			return err
		}
	}
	return nil
}

 // Decoding the Block
 func (b *Block) Decode (dec Decoder[*Block]) error{
	return dec.Decode(b)
 }

// Encoding the Block
func (b *Block) Encode (enc Encoder[*Block]) error{
	return enc.Encode(b)
 }

 func (b *Block) Hash(hasher Hasher[*Header]) types.Hash{
	if b.hash.IsZero(){
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}







