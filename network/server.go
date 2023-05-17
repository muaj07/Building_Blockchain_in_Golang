package network
import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}
type Server struct {
	ServerOpts
	rpcCh chan RPC
	quitCh chan struct{}

}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)
	free:
	for {
		select {
		case rpc := <- s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <- ticker.C:
			fmt.Println("Do stuff every x seconds")
		case <- s.quitCh:
			break free
		}
	}
	fmt.Println("Server Shutdown")
	
}

func (s *Server) initTransports () {
	for _, tr := range s.Transports {
		go func (tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
				}
			}(tr)
			}
		}