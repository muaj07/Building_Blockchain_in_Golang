package network
import(
	"net"
	"fmt"
	"bytes"
)


type TCPPeer struct {
	conn net.Conn
	Outgoing bool
}
func (p *TCPPeer) Send(b []byte) error{
	_,err := p.conn.Write(b)
	return err	
}

func (p *TCPPeer) readLoop(rpcCh chan RPC){
	buf := make([]byte, 2048)
	for {
		n, err := p.conn.Read(buf)
		if err != nil{
			fmt.Printf("Read error : %s", err)
			continue
		}
		msg := buf[:n]
		rpcCh <- RPC {
			From: p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
		//Handle the message Here
	}
}


type TCPTransport struct {
	peerCh chan *TCPPeer
	listenAddr string
	listener net.Listener
}

func NewTCPTransport (addr string, peerCh chan *TCPPeer) *TCPTransport{
	return &TCPTransport{
		peerCh: peerCh,
		listenAddr: addr,
	}
}


//This is a method called "Start" that belongs to a struct 
//called "TCPTransport". It starts a TCP server by listening 
//on a specified address and returns an error if there is one. 
//If there is no error, it sets the listener to the returned 
//listener, prints a message to the console indicating that 
//it is listening on the specified address, and then starts 
//an accept loop in a new goroutine.
func (t *TCPTransport) Start() error{
	ln, err := net.Listen("tcp", t.listenAddr)
	if err!= nil{
		return err
	}
	//assign the TCP Server "net.Listen" method "Listener" to the "t.listener"
	t.listener = ln
	fmt.Println("TCP Transport listening to Port:", t.listenAddr)
	//Accept the connection
	go t.acceptLoop()
	return nil
}

// acceptLoop listens for incoming TCP connections and creates a new TCPPeer for each connection.
// It runs indefinitely until an error occurs.
func (t *TCPTransport) acceptLoop() {
    for {
        // Wait for a new connection
        conn, err := t.listener.Accept()
        if err != nil {
            // If there was an error accepting the connection, log it and continue accepting connections
            fmt.Printf("Accept error from %+v\n", conn)
            continue
        }

        // Create a new TCPPeer for the connection
        peer := &TCPPeer{
            conn: conn,
        }

        // Add the new peer to the peer channel so it can be processed by the transport
        t.peerCh <- peer

        // Log the new incoming connection
        fmt.Printf("New Incoming TCP Connection %+v\n", conn)

        
    }
}



// func (t *TCPTransport) readLoop(peer *TCPPeer){
// 	buf := make([]byte, 2048)
// 	for {
// 		n, err := peer.conn.Read(buf)
// 		if err != nil{
// 			fmt.Printf("Read error : %s", err)
// 			continue
// 		}
// 		msg := buf[:n]
// 		fmt.Println(string(msg))
// 		//Handle the message Here
// 	}

// }
