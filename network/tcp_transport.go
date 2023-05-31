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

// readLoop listens for incoming data on the TCPPeer's connection, reads it into a buffer, and sends it to a channel for processing.
func (p *TCPPeer) readLoop(rpcCh chan RPC) {
    // create a buffer to hold incoming data
    buf := make([]byte, 2048)

    // loop indefinitely
    for {
        // read data from the connection into the buffer
        n, err := p.conn.Read(buf)
        // if there is an error reading, print a message and continue to next iteration
        if err != nil {
            fmt.Printf("Read Error : %s", err)
            continue
        }

        // extract the actual message from the buffer
        msg := buf[:n]

        // create an RPC message with the source address and message payload, and send it to the designated channel for processing
        rpcCh <- RPC {
            From: p.conn.RemoteAddr(),
            Payload: bytes.NewReader(msg),
        }

        // message is now in the channel and can be processed
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


// Start is a method that starts a TCP server by listening on a specified address.
// It belongs to a struct called "TCPTransport". It returns an error if there is one.
// If there is no error, it sets the listener to the returned listener, prints a 
//message to the console indicating that it is listening on the specified address, 
//and then starts an accept loop in a new goroutine.
func (t *TCPTransport) Start() error {
    // Listen for incoming TCP connections on the specified address
    ln, err := net.Listen("tcp", t.listenAddr)
    if err != nil {
        // Return an error if listening fails
        return err
    }
    // Set the listener to the returned listener
    t.listener = ln
    // Print a message indicating that the server is listening on the specified address
    fmt.Println("TCP Transport listening to Port:", t.listenAddr)
    // Start an accept loop in a new goroutine
    go t.acceptLoop()
    // Return nil since there was no error
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


