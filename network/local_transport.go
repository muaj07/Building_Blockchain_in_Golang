package network

import (
	"sync"
	"fmt"
	"bytes"
    "net"
)

type LocalTransport struct {
	addr net.Addr
	ConsumeCh chan RPC
	lock sync.RWMutex
	peers map[net.Addr]*LocalTransport
}

// NewLocalTransport creates a new instance of the LocalTransport struct
// with the specified net.Addr and initializes its fields
func NewLocalTransport(addr net.Addr) *LocalTransport {
    // Initialize a new LocalTransport struct with the specified net.Addr
    // and default values for its fields
    return &LocalTransport {
        addr:      addr,
        ConsumeCh: make(chan RPC, 1024),
        peers:     make(map[net.Addr]*LocalTransport),
    }
}

func(t *LocalTransport) Broadcast(payload []byte) error{
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err!=nil {
			return err
		}
	}
	return nil
}
// Consume returns a channel that can be used to receive RPCs from the local transport.
func (t *LocalTransport) Consume() <- chan RPC {
	return t.ConsumeCh
}
// Addr returns the network address of the local transport.
func (t *LocalTransport) Addr() net.Addr {
	return t.addr
}

// SendMessage sends an RPC message to the specified peer.
// The `to` parameter is the address of the peer.
// The `payload` parameter is the Data to send.
func (t *LocalTransport) SendMessage(to net.Addr, payload []byte) error {
    // Lock the peers map to prevent concurrent access.
    t.lock.RLock()
    defer t.lock.RUnlock()
    if t.addr==to{
        return nil
    }

    // Attempt to find the peer in the map.
    peer, ok := t.peers[to]
    if !ok {
        return fmt.Errorf("Error! %s could not send Msg to unknown Peer: %s", t.addr, to)
    }

    // Send the message via the peer's ConsumeCh channel
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
