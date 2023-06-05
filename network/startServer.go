package network

import(
	"fmt"
	"time"
	//"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
	//"github.com/go-kit/log"
	//"bytes"
	//"os"
	"net"
)

// Start starts the server.
// It initializes the transports, starts the ticker, and listens for incoming RPC requests.
// If the server is a validator, it creates a new block for each tick of the ticker.
func (s *Server) Start() {
    // Initialize TCP Transport
    s.TCPTransport.Start() 
	time.Sleep(1 * time.Second)
	s.bootStrapNetwork()

	s.Logger.Log(
		"msg", "Accepting TCP connection ",
		//"Server ID", s.ID,
		"Listening Address", s.ListenAddr,
	)

    // free is the name for the for loop
free:
    for {
        select {
		case peer := <-s.peerCh:
			s.peerMap[peer.conn.RemoteAddr()] = peer
			go peer.readLoop(s.rpcCh)
			s.Logger.Log(
				"peer", peer,
			)
			if err := s.SendGetStatusMessage(peer); err!=nil{
				s.Logger.Log(
									"msg", "Failed to Send GetStatusMessage",
									"err", err,
							)
				continue			
			}
			s.Logger.Log(
				"msg", "Peer Added to the Server",
				"outgoing", peer.Outgoing,
				"addr", peer.conn.RemoteAddr(),
			)
        // Handle incoming RPC requests
        case rpc := <-s.rpcCh:
			//the msg is decoded and contain the Transaction struct, which is
			//defined in "transaction.go" file of the Core Package
			msg, err := s.RPCDecodeFunc(rpc)
            if err != nil {
                s.Logger.Log(
					"Error", err,
				)
				continue
            }
			// "RPCProcessor" is an interface with single method "ProcessMessage"
			// defined in the "rpc.go" file.
			// The server implement the "RPCProcessor" interface below
			// in this file
			if err := s.RPCProcessor.ProcessMessage(msg); err!=nil {
				if err!= core.ErrBlockKnown{
				s.Logger.Log(
					"Error", err,
				)
				}
			}
        // Quit gracefully if quitCh is closed
        case <-s.quitCh:
            break free // break the free for loop
        }
    }
    // Print message after server is shut down
	s.Logger.Log(
		"msg", "Server Shutting down",
	)
}


func (s *Server) bootStrapNetwork(){
	for _, addr := range s.SeedNodes{
		fmt.Println("BootstrapNetwork trying to connect to Seed Nodes", addr)
		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err!=nil{
				fmt.Printf("Could not connect to %+v\n", conn)
				return
			}
			s.Logger.Log(
                "msg", "Connected to Seed Node",
                "addr", addr,
            )
			s.peerCh <- &TCPPeer{
				conn: conn,
			}
			fmt.Printf("Bootstrap %+v\n", conn)
		}(addr)
		
	}
}