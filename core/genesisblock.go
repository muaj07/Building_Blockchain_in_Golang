package core

import (
	"github.com/muaj07/transport/types"
    //"fmt"
)
// GenesisBlock creates and returns the first block of the blockchain.
func GenesisBlock() *Block {
    //fmt.Println("Creating genesis block...")
    // Create the header for the block.
    header := &Header{
        Version:  1,             // The version number of the blockchain.
        DataHash: types.Hash{},  // The hash of the block's data.
        Height:   0,             // The height of the block in the blockchain.
        TimeStamp: 00000,       // The Unix timestamp of when the block was created.
    }
    // Create the block using the header.
    //fmt.Println("Creating block...")
    // No tx included in the NewBlock since this is a Genesis block
    b, _ := NewBlock(header, nil)

    // Return the block.
    //fmt.Println("Genesis block created.")
    return b
}
