package network

type NetAddr string

type Transport interface {
	Connect(Transport) error
	Consume()	<- chan RPC
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}