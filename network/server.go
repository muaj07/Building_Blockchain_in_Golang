package network
import (
	//"fmt"
	"time"
	"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
	"github.com/go-kit/log"
	"github.com/muaj07/transport/types"
	"bytes"
	"os"
)

var defaultBlockTime = 4*time.Second

// This code defines a struct called ServerOpts that contains 
//fields for an RPC handler, a slice of transports, a block time 
// duration, and a private key for cryptography
type ServerOpts struct {
	ID string
	Logger log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor RPCProcessor
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
	memPool *TxPool
	chain *core.Blockchain
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

// NewServer returns a new instance of Server with the provided options.
// If block time is not set, use default block time.
// Set default RPC handler if none is provided.
func NewServer(opts ServerOpts) (*Server, error) {
	if opts.Logger == nil{
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
    // If block time is not set, use default block time.
    if opts.BlockTime == time.Duration(0) {
        opts.BlockTime = defaultBlockTime
    }
	// Set RPCDecodeFunc if none is provided.
    if opts.RPCDecodeFunc == nil {
        opts.RPCDecodeFunc = DefaultRPCDecodeFunc
    }
	chain, err := core.NewBlockchain(genesisBlock())
	if err != nil {
		return nil, err
	}
    // Create a new server instance.
    s := &Server{
        ServerOpts: opts,
		chain : chain,
		//Initial size of the memPool set to 100 txs
        memPool:    NewTxPool(100),
        // Validator needs privatekey to sign the blocks
        isValidator: opts.PrivateKey != nil,
        rpcCh:       make(chan RPC),
        quitCh:      make(chan struct{}, 1),
    }

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

// Start starts the server.
// It initializes the transports, starts the ticker, and listens for incoming RPC requests.
// If the server is a validator, it creates a new block for each tick of the ticker.
func (s *Server) Start() {
    // Initialize transports
    s.initTransports()

    // Main loop
free:
    for {
        select {
        // Handle incoming RPC requests
        case rpc := <-s.rpcCh:
			//the msg is decoded and contain the Transaction struct, which is
			//defined in "transaction.go" file of the Core Package
			msg, err := s.RPCDecodeFunc(rpc)
            if err != nil {
                s.Logger.Log(
					"Error", err,
				)
            }
			// "RPCProcessor" is an interface with single method "ProcessMessage"
			// defined in the "rpc.go" file.
			// The server implement the "RPCProcessor" interface below
			// in this file
			if err := s.RPCProcessor.ProcessMessage(msg); err!=nil {
				s.Logger.Log(
					"Error", err,
				)
			}
        // Quit gracefully if quitCh is closed
        case <-s.quitCh:
            break free
        }
    }
    // Print message after server is shut down
	s.Logger.Log(
		"msg", "Server Shutting down",
	)
}

func (s *Server) validatorLoop(){
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log(
		"msg", "Starting Validator ",
		"blockTime", s.BlockTime,
	)
	for{
		<-ticker.C
		s.createNewBlock()
	}
}


func (s *Server) createNewBlock() error{
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}
	// For now, all the available txs in a mempool are
	// included in the block, which will be update in the future
	tx := s.memPool.Pending()
	block, err := core.NewBlockFromPrevHeader(currentHeader, tx)
	if err != nil {
		return nil 
	}
	if err := block.Sign(*s.PrivateKey); err!=nil{
		return err
	}
	if err := s.chain.AddBlock(block); err != nil{
		return err
	}
	s.memPool.ClearPending()
	go s.broadcastBlock(block)
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


// ProcessMessage processes a message received by the server
// and returns an error (if any).
func (s *Server) ProcessMessage(msg *DecodeMessage) error {
    // Determine whether the message is a transaction or a block.
    switch tx := msg.Data.(type) {
    case *core.Transaction:
        // Process the transaction.
        return s.processTransaction(tx)
    case *core.Block:
        // Process the block.
        return s.processBlock(tx)
    }
    // The message was not a transaction or a block, so return nil.
    return nil
}


func (s *Server) processBlock (b *core.Block) error{
	if err := s.chain.AddBlock(b); err != nil{
		return err
	}
	go s.broadcastBlock(b)
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error{
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Contains(hash){
		return nil
	}
	if err := tx.Verify(); err!=nil{
		return err
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	//s.Logger.Log(
	//	"msg", "Adding new tx to mempool", 
	//	"hash", hash, 
	//	"mempoolLength", s.memPool.PendingCount(),
	//)
	go s.broadcastTx(tx)
	s.memPool.Add(tx)
	return nil
}


// broadcastBlock broadcasts a block to all connected peers.
// Returns an error if there was an issue encoding the block 
//or broadcasting the message.
func (s *Server) broadcastBlock(b *core.Block) error {
    buf := &bytes.Buffer{}
    if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
        return err
    }
    msg := NewMessage(MessageTypeBlock, buf.Bytes())
    return s.broadcast(msg.Bytes())
}


// broadcastTx encodes a transaction and broadcasts it to all connected peers.
// If an error occurs, it is returned.
func (s *Server) broadcastTx(tx *core.Transaction) error {
    // Create a buffer to hold the encoded transaction.
    buf := &bytes.Buffer{}

    // Encode the transaction using the GobTxEncoder.
    if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
        return err
    }

    // Create a new message with the encoded transaction.
    msg := NewMessage(MessageTypeTx, buf.Bytes())

    // Broadcast the message to all connected peers.
    return s.broadcast(msg.Bytes())
}


// broadcast sends Payload to all connected transports in Server s
func (s *Server) broadcast(Payload []byte) error {
    // iterate over each transport in the server and broadcast the payload
    for _, tr := range s.Transports {
        if err := tr.Broadcast(Payload); err != nil {
            return err
        }
    }
    // no errors occurred, return nil
    return nil
}


// genesisBlock creates and returns the first block of the blockchain.
func genesisBlock() *core.Block {
    // Create the header for the block.
    header := &core.Header{
        Version:  1,             // The version number of the blockchain.
        DataHash: types.Hash{},  // The hash of the block's data.
        Height:   0,             // The height of the block in the blockchain.
        TimeStamp: 00000,       // The Unix timestamp of when the block was created.
    }
    // Create the block using the header.
	// No tx included in the NewBlock since this is a Genesis block
    b, _ := core.NewBlock(header, nil)

    // Return the block.
    return b
}

