package core
import (
	"fmt"
)

type Validator interface {
	// ValidateBlock checks if a given block is valid
	ValidateBlock(*Block) error
}

type BlockValidator struct{
	bc *Blockchain
}

// NewBlockValidator returns a new instance of BlockValidator with a reference to the provided Blockchain.
//
// Parameters:
// - bc (*Blockchain): a pointer to the Blockchain instance to be used for validation.
// Returns:
// - (*BlockValidator): a pointer to the newly created BlockValidator instance.
func NewBlockValidator (bc *Blockchain) *BlockValidator{
	return &BlockValidator {
		bc: bc,
	}
}


// ValidateBlock validates the given Block using the BlockValidator.
// It checks if the block is already present in the chain and verifies the block's signature.
//
// Parameters:
// - bv (*BlockValidator): a pointer to a BlockValidator object.
// - b (*Block): a pointer to a Block object to validate.
//
// Returns:
// - error: an error that occurred during validation, or nil if validation succeeded.
func (bv *BlockValidator) ValidateBlock(b *Block) error {
    // Check if the block is already present in the chain.

    if bv.bc.HasBlock(b.Height) {
        return fmt.Errorf("Chain already contains Block#%d with Hash ==> %s", b.Height, b.Hash(BlockHasher{}))
    }
	if b.Height != bv.bc.Height()+1{
		return fmt.Errorf("Block ==> %s with Height # %d is too high, Current height ==>%d)", b.Hash(BlockHasher{}), b.Height, bv.bc.Height()+1)
	}
	prevHeader, err := bv.bc.GetHeader(b.Height-1)
	if err != nil {
		return err
	}
	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash{
		return fmt.Errorf("The hash %s of Previous Block is INVALID", b.PrevBlockHash)
	}
    // Verify the block's signature.
    if err := b.Verify(); err!= nil{
        return err
    }
	fmt.Printf("Done Basic Verification for Block#%d\n", b.Height)
    return nil
}
