package network

import (
	"sync"
	"fmt"
)

type LocalTransport struct {
	addr NetAddr
	ConsumeCh chan	RPC
	lock sync.RWMutex
	peers map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr: addr,
		ConsumeCh: make(chan RPC, 1024),
		peers: make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume () <- chan RPC{
	return t.ConsumeCh
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	peer, ok:= t.peers[to]
	if !ok {
		return fmt.Errorf("Error! %s could not send Msg to %s", t.addr,to)
	}
	peer.ConsumeCh <- RPC {
		From: t.addr,
		Payload: payload,
	}
	return nil
}

func (t *LocalTransport) Connect (tr Transport) error {
	trans := tr.(*LocalTransport)	
	t.lock.Lock()
	defer t.lock.Unlock()
	t.peers[tr.Addr()]= trans
	return nil
}