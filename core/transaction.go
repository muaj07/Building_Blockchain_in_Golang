package core
import(
	"io"
	"github.com/muaj07/transport/crypto"
	"fmt"
	"github.com/muaj07/transport/types"
)

type Transaction struct {
	Data []byte
	From crypto.PublicKey
	Signature *crypto.Signature
	//Cached version of the tx data Hash
	// We don't want them public
	hash types.Hash
	firstSeen int64
}


// NewTransaction creates a new Transaction instance with the given data.
// The data parameter should be a byte array representing the transaction data
func NewTransaction (data []byte) *Transaction{
	return &Transaction{
		Data: data,
	}
}


// Hash returns the hash of the transaction data
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero(){
		tx.hash = hasher.Hash(tx)
	}
    return tx.hash
}



// Sign the Transaction using the "Sign" method from
// the "keypair.go" file in the crypto Folder
// The "Sign method" returns signture (or error), which is assigned to sig
// and then stored in "tx.signature", while the Publickey for the given/used
// private key is stored in "tx.publickey"

func (tx *Transaction) Sign(privkey crypto.PrivateKey) error{
	sig, err := privkey.Sign(tx.Data)
	if err!=nil{
		return err
	}
	tx.Signature = sig
	tx.From = privkey.PublicKey()
	return nil
}

// Verify the trx by checking
// if the "transaction signature" exist and then
// if the transaction is signed by the correct private key

func (tx *Transaction) Verify() error {
	if tx.Signature == nil{
		return fmt.Errorf ("Transaction signature is Nil")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf ("Invalid transaction signature")
	}
	return nil
}

func (tx *Transaction) EncodeBinary (w io.Writer) error {
	return nil
	
}

func (tx *Transaction) DecodeBinary (r io.Reader) error {
	return nil
}


func (tx *Transaction) Decode(dec Decoder[*Transaction]) error{
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error{
	return enc.Encode(tx)
}
func (tx *Transaction) SetFirstSeen (t int64) {
	tx.firstSeen = t
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}