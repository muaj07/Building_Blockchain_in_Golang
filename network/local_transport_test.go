package network

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t,tra.peers[trb.addr], trb)
	assert.Equal(t,trb.peers[tra.addr], tra)
}

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("Hello Enoda!" + "\n")
	assert.Nil(t, tra.SendMessage(trb.addr, msg))
	rpc := <- trb.Consume()
	assert.Equal(t, msg, rpc.Payload)
	assert.Equal(t, tra.addr, rpc.From)	
}