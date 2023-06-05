package core

import (
//"encoding/binary"
"github.com/muaj07/transport/types"
//"bytes"
//"crypto/sha256"
"fmt"
"github.com/muaj07/transport/crypto"
//"encoding/gob"
//"time"
)

type Block struct {
	*Header
	Transactions []*Transaction
	Validator crypto.PublicKey //Public key of the Block validator
	Signature *crypto.Signature //Block signature
	// cached version of the Header hash
	hash types.Hash
}

//Constructor function for Block struct
// Only two fields are included in the signature of the constructor function
// other values will be set/updated later in the code

func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header: h,
		Transactions: txx,
	}, nil	
}




func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

// Sign the Block using the "Sign" method from
// the "keypair.go" file in the crypto Folder.
// The "Sign method" returns (*Signature, error), 
// which is assigned to sig and then stored in 
// "b.signature", while the Public key of the 
// Validator is stored in "b.Validator"

func (b *Block) Sign(privkey crypto.PrivateKey) error{
	sig, err := privkey.Sign(b.Header.Bytes())
	if err!=nil{
		return err
	}
	// set the validator and signature fields of the block struct
	b.Signature = sig
	b.Validator = privkey.PublicKey()
	return nil
}


// Verify the Block by checking
// if the "Validator signature" exist and then
// if the transaction is signed by the correct (i.e. of the Validator) Private key

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("Block has no Signature")
	}
	//verify takes the public key and data, create
	//new signature and compare with "b.Signature"
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf ("Invalid Block signature")
	}
	// Verify all the txs in a block
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err!=nil {
			return err
		}
	}
	//verify the data Hash of the txs if it matches with the calculated hash
	dataHash, err := CalculateDataHash(b.Transactions)
	if err!= nil {
		return err
	}
	if dataHash != b.Header.DataHash {
		return fmt.Errorf("Block (%s) has invalid data hash", b.Hash(BlockHasher{}))
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





