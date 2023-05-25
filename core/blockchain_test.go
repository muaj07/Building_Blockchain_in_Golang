package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/muaj07/transport/types"
	//"github.com/go-kit/log"
	//"github.com/muaj07/transport/crypto"
)



//This code snippet defines a function called newBlockchainWithGenesis 
//that takes a *testing.T pointer as an argument and returns a *Blockchain pointer. 
//The function creates a new blockchain by calling NewBlockchain with a randomly 
//generated block of height 0 as its argument. Finally, the function returns the 
//newly created blockchain. If there are any errors during the creation process, 
//the function will fail the test with an assertion error.

// newBlockchainWithGenesis creates a new blockchain with a genesis block for testing purposes.
// It takes a testing object as an argument and returns a pointer to the new blockchain.
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
    // Create a new blockchain with a random genesis block.
    bc, err := NewBlockchain(randomBlock(t, 0, types.Hash{}))
    // If there was an error, fail the test.
    assert.Nil(t, err)
    // Return the new blockchain.
    return bc
}

// getPrevBlockHash returns the hash of the previous block at the given height
// It takes in a testing.T object, a Blockchain pointer, and a height of the block
func getPrevBlockHash (t *testing.T, bc *Blockchain, height uint32) types.Hash{
	prevHeader, err := bc.GetHeader(height-1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}

// TestNewBlockchain tests the creation of a new blockchain
func TestNewBlockchain(t *testing.T) {
    // Create a new blockchain with genesis block
    bc := newBlockchainWithGenesis(t)

    // Check that the validator is not nil
    assert.NotNil(t, bc.validator)

    // Check that the height of the blockchain is 0
    assert.Equal(t, bc.Height(), uint32(0))
}

func TestAddBlock(t *testing.T){
	 // Create a new blockchain with genesis block
	 bc := newBlockchainWithGenesis(t)
	 lenBlocks := 100
	 for i:=0; i<lenBlocks; i++{
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
	 	assert.Nil(t, bc.AddBlock(block))
	 }
	 assert.Equal(t, len(bc.headers), lenBlocks+1)
	 assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

// TestHasBlock tests the HasBlock method of the Blockchain struct
func TestHasBlock(t *testing.T) {
    // create a new blockchain with genesis block
    bc := newBlockchainWithGenesis(t)

    // test that genesis block is present
    assert.True(t, bc.HasBlock(0))

    // test that non-existent block is not present
    assert.False(t, bc.HasBlock(1))
}

func TestAddBlockTooHigh (t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}