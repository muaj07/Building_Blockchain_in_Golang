package network

import(
	"fmt"
	"github.com/muaj07/transport/core"
	//"github.com/go-kit/log"
	"bytes"
	//"os"
	//"net"
)

// BroadcastTx encodes a transaction and broadcasts it to all connected peers.
// If an error occurs, it is returned.
func (s *Server) BroadcastTx(tx *core.Transaction) error {
    // Create a buffer to hold the encoded transaction.
    buf := &bytes.Buffer{}

    // Encode the transaction using the GobTxEncoder.
    if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
        return err
    }

    // Create a new message with the encoded transaction.
    msg := NewMessage(MessageTypeTx, buf.Bytes())

    // Broadcast the message to all connected peers.
    return s.Broadcast(msg.Bytes())
}


// BroadcastBlock broadcasts a block to all connected peers.
// Returns an error if there was an issue encoding the block 
//or broadcasting the message.
func (s *Server) BroadcastBlock(b *core.Block) error {
    buf := &bytes.Buffer{}
    if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
        return err
    }
    msg := NewMessage(MessageTypeBlock, buf.Bytes())
    return s.Broadcast(msg.Bytes())
}


// Broadcast sends Payload to all connected transports in Server s
func (s *Server) Broadcast(Payload []byte) error {
    //iterate over each transport in the peerMap
	s.mu.RLock()
	defer s.mu.RUnlock()
    for netAddr, peer := range s.peerMap {
        if err := peer.Send(Payload); err != nil {
            fmt.Printf("Peer send error ==> addr %s [%s]\n",netAddr,err)
        }
    }
    return nil
}