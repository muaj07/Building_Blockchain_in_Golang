package network
import (
	//"fmt"
	"time"
	"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
	"github.com/go-kit/log"
	//"github.com/muaj07/transport/types"
	//"encoding/gob"
	//"bytes"
	"os"
	"net"
	"sync"
)

var defaultBlockTime = 4*time.Second


type ServerOpts struct {
	SeedNodes []string
	ListenAddr string
	ID string
	TCPTransport	*TCPTransport
	Logger log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
	BlockTime	time.Duration
	PrivateKey *crypto.PrivateKey
}


type Server struct {
	TCPTransport *TCPTransport
	peerCh chan *TCPPeer
	mu sync.RWMutex
	peerMap map[net.Addr]*TCPPeer
	ServerOpts
	memPool *TxPool
	chain *core.Blockchain
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}


func NewServer(opts ServerOpts) (*Server, error) {
	if opts.Logger == nil{
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "Node_Addr", opts.ID)
	}
    // If block time is not set, use default block time.
    if opts.BlockTime == time.Duration(0) {
        opts.BlockTime = defaultBlockTime
    }
	// Set RPCDecodeFunc if none is provided.
    if opts.RPCDecodeFunc == nil {
        opts.RPCDecodeFunc = DefaultRPCDecodeFunc
    }
	
	chain, err := core.NewBlockchain(core.GenesisBlock())
	if err != nil {
		return nil, err
	}
	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

    // Create a new server instance.
    s := &Server{
		TCPTransport: tr,
		peerCh: peerCh,
		peerMap:	make(map[net.Addr]*TCPPeer),
        ServerOpts: opts,
		chain : chain,
        memPool:     NewTxPool(100), //Initial size of the memPool set to 100 txs
        isValidator: opts.PrivateKey != nil, //Validator needs privatekey to sign the blocks
        rpcCh:       make(chan RPC),
        quitCh:      make(chan struct{}),
    }
	s.TCPTransport.peerCh = peerCh
    // if there is no processor assigned in the Server options
	// then use the server as a default processor
	if s.RPCProcessor == nil{
		s.RPCProcessor= s
		}
	if s.isValidator{
		go s.validatorLoop()
	}
    // Return the server instance.
    return s, nil
}

// func (s *Server) bootStrapNodes(){
// 	for _, tr:= range s.Transports {
// 		if s.Transport.Addr()!= tr.Addr(){
// 			if err := s.Transport.Connect(tr); err!=nil {
// 			s.Logger.Log(
// 				"Error", "Could not Connect to Remote",
// 				"err", err,
// 			)
// 			}
// 			s.Logger.Log(
// 				"msg", "Connect to Remote",
// 				"host addr", s.Transport.Addr(),
// 				"Remote addr", tr.Addr(),
// 			)
// 			// send the getStatusMessage so we can sync (if needed)
// 			if err := s.SendGetStatusMessage(tr); err != nil{
// 				s.Logger.Log(
// 					"Error", err,
// 				)
// 			}
// 		}
// 	}	
// }

// initTransports initializes the transports for the server, 
// starting a goroutine for each transport to consume RPCs.
// func (s *Server) initTransports() {
//     // Loop through each transport in the server's Transports slice.
//     for _, tr := range s.Transports {
//         // Start a goroutine for each transport.
//         go func(tr Transport) {
//             // Loop through each RPC consumed by the transport.
//             for rpc := range tr.Consume() {
//                 // Send the RPC to the server's rpcCh channel.
//                 s.rpcCh <- rpc
//             }
//         }(tr)
//     }
// }


