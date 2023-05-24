package core

import (
	"fmt"
	"github.com/go-kit/log"
	"sync"
	
)

type Blockchain struct{
	logger log.Logger
	store Storage
	lock  sync.RWMutex
	headers []*Header
	validator Validator
}

// NewBlockchain creates a new instance of the Blockchain struct and returns a pointer to it.
//
// Returns:
// *Blockchain: A pointer to the newly created Blockchain instance.
func NewBlockchain (l log.Logger, genesis *Block) (*Blockchain, error){
	bc:= &Blockchain{
		logger: l,
		store: NewMemoryStorage(),
		headers: []*Header{},
		//Separate the "initialization" of the validator field and set it 
		//after the bc instance is fully initialized.
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.AddBlockWithoutValidation(genesis)
	return bc, err
}

// SetValidator sets the validator for the Blockchain.
// v is the Validator to be set for the Blockchain.
func (bc *Blockchain) SetValidator (v Validator) {
	bc.validator = v
}


// AddBlock adds a validated block to the blockchain.
func (bc *Blockchain) AddBlock(b *Block) error {
    // Validate the block before adding it
    if err := bc.validator.ValidateBlock(b); err != nil {
        return err
    }
    // Add the block to the blockchain.
    return bc.AddBlockWithoutValidation(b)
}


func (bc *Blockchain) GetHeader (height uint32) (*Header, error) {
	if height > bc.Height(){
		return nil, fmt.Errorf("Given height (%d) is too high", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.headers[height], nil
}


// AddBlockWithoutValidation adds a block to the blockchain without validating it.
// It appends the block's header to the blockchain's headers slice and stores the block.
func (bc *Blockchain) AddBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
    bc.headers = append(bc.headers, b.Header) // append the block's header to the headers slice
	bc.lock.Unlock()
	bc.logger.Log (
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions), 
	)
    return bc.store.Put(b) // store the block in the blockchain's store
}

// Height returns the current height of the blockchain.
// This function does not take any parameters.
// It returns a uint32 value representing the height of the blockchain.
func (bc *Blockchain) Height() uint32{
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers)-1)
}

// HasBlock returns a boolean indicating whether a block exists at the given height in the blockchain.
//
// height: The height of the block to check.
// Returns: A boolean value indicating whether a block exists at the given height.
func (bc *Blockchain) HasBlock(height uint32) bool{
	return height <= bc.Height()
}