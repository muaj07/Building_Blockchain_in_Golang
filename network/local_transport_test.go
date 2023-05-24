package network

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	b, err := ioutil.ReadAll(rpc.Payload)
	
	assert.Nil(t,err)
	assert.Equal(t, msg, b)
	assert.Equal(t, rpc.From, tra.addr)	
}
func TestBroadcast(t *testing.T){
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")
	tra.Connect(trb)
	tra.Connect(trc)
	trb.Connect(tra)
	trb.Connect(trc)
	msg := []byte("Hello Enoda!" + "\n")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <- trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t,err)
	assert.Equal(t, b, msg)
	assert.Equal(t, len(tra.peers), 2)	
	assert.Equal(t, len(trb.peers), 2)	

	rpcc := <- trc.Consume()
	c, err := ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t,err)
	assert.Equal(t, c, msg)
	assert.Equal(t, len(tra.peers), 2)	
	assert.Equal(t, len(trb.peers), 2)
}