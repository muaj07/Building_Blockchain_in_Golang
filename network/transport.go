package network
import "net"

//type NetAddr string

type Transport interface {
	Connect(Transport) error
	Consume()	<- chan RPC
	SendMessage(net.Addr, []byte) error
	Broadcast([]byte) error
	Addr() net.Addr
}