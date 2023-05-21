package network

import (
	"sync"
	"fmt"
	//"io"
	"bytes"
)

type LocalTransport struct {
	addr NetAddr
	ConsumeCh chan RPC
	lock sync.RWMutex
	peers map[NetAddr]*LocalTransport
}

// NewLocalTransport creates a new instance of the LocalTransport struct
// with the specified NetAddr and initializes its fields
func NewLocalTransport(addr NetAddr) *LocalTransport {
    // Initialize a new LocalTransport struct with the specified NetAddr
    // and default values for its fields
    return &LocalTransport {
        addr:      addr,
        ConsumeCh: make(chan RPC, 1024),
        peers:     make(map[NetAddr]*LocalTransport),
    }
}

// Consume returns a channel that can be used to receive RPCs from the local transport.
func (t *LocalTransport) Consume() <- chan RPC {
	return t.ConsumeCh
}
// Addr returns the network address of the local transport.
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

// SendMessage sends an RPC message to the specified peer.
// The `to` parameter is the address of the peer.
// The `payload` parameter is the Data to send.
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
    // Lock the peers map to prevent concurrent access.
    t.lock.RLock()
    defer t.lock.RUnlock()

    // Attempt to find the peer in the map.
    peer, ok := t.peers[to]
    if !ok {
        return fmt.Errorf("Error! %s could not send Msg to %s", t.addr, to)
    }

    // Send the message via the peer's ConsumeCh channel.
    peer.ConsumeCh <- RPC {
        From: t.addr,
        Payload: bytes.NewReader(payload),
    }
    // Return nil to indicate success.
    return nil
}
// Connect adds the given transport to the local transport's peer list.
// It returns an error if the transport is already connected.
func (t *LocalTransport) Connect(tr Transport) error {
    // Cast the transport to a LocalTransport
    trans := tr.(*LocalTransport)

    t.lock.Lock()
    defer t.lock.Unlock()

    // Add the transport to the peer list
    t.peers[tr.Addr()] = trans
    return nil
}
