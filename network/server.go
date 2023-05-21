package network
import (
	"fmt"
	"time"
	"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5*time.Second

// This code defines a struct called ServerOpts that contains 
//fields for an RPC handler, a slice of transports, a block time 
// duration, and a private key for cryptography
type ServerOpts struct {
	RPCHandler RPCHandler
	Transports []Transport
	BlockTime	time.Duration
	PrivateKey *crypto.PrivateKey
}

// The Server struct contains some fields, including ServerOpts
// (a field of another struct type), a blockTime duration, a
// pointer to a TxPool struct, a boolean field, and two channels
// (one for receiving RPC messages and another for quitting).

type Server struct {
	ServerOpts
	blockTime	time.Duration
	memPool *TxPool
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}

}

// NewServer returns a new instance of Server with the provided options.
// If block time is not set, use default block time.
// Set default RPC handler if none is provided.
func NewServer(opts ServerOpts) *Server {
    // If block time is not set, use default block time.
    if opts.BlockTime == time.Duration(0) {
        opts.BlockTime = defaultBlockTime
    }

    // Create a new server instance.
    s := &Server{
        ServerOpts: opts,
        blockTime:  opts.BlockTime,
        memPool:    NewTxPool(),
        // Validator needs privatekey to sign the blocks
        isValidator: opts.PrivateKey != nil,
        rpcCh:       make(chan RPC),
        quitCh:      make(chan struct{}, 1),
    }

    // Set default RPC handler if none is provided.
    if opts.RPCHandler == nil {
        opts.RPCHandler = NewDefaultRPCHandler(s)
    }
	s.ServerOpts= opts
    // Return the server instance.
    return s
}

// Start starts the server.
// It initializes the transports, starts the ticker, and listens for incoming RPC requests.
// If the server is a validator, it creates a new block for each tick of the ticker.
func (s *Server) Start() {
    // Initialize transports
    s.initTransports()

    // Start ticker
    ticker := time.NewTicker(s.BlockTime)

    // Main loop
free:
    for {
        select {
        // Handle incoming RPC requests
        case rpc := <-s.rpcCh:
            if err := s.RPCHandler.HandleRPC(rpc); err != nil {
                logrus.Error(err)
            }
        // Create new block for each tick of the ticker
        case <-ticker.C:
            if s.isValidator {
                s.createNewBlock()
            }
        // Quit gracefully if quitCh is closed
        case <-s.quitCh:
            break free
        }
    }
    // Print message after server is shut down
    fmt.Println("Server Shutdown")
}


func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error{
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash){
		logrus.WithFields(logrus.Fields{
			"hash": hash,
			"mempool length": s.memPool.Len(),
		}).Info("Transaction Already in Mempool")
		return nil
	}
	if err := tx.Verify(); err!=nil{
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	logrus.WithFields(logrus.Fields{
		"hash": hash,
		"mempool length": s.memPool.Len(),
	}).Info("Adding new tx to the mempool")
	return s.memPool.AddTx(tx)
}


func (s *Server) createNewBlock() error{
	logrus.Info("Create new Block")
	return nil
}


// initTransports initializes the transports for the server, 
// starting a goroutine for each transport to consume RPCs.
func (s *Server) initTransports() {
    // Loop through each transport in the server's Transports slice.
    for _, tr := range s.Transports {
        // Start a goroutine for each transport.
        go func(tr Transport) {
            // Loop through each RPC consumed by the transport.
            for rpc := range tr.Consume() {
                // Send the RPC to the server's rpcCh channel.
                s.rpcCh <- rpc
            }
        }(tr)
    }
}
