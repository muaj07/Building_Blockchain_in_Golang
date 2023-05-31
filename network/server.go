package network
import (
	"fmt"
	"time"
	"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
	"github.com/go-kit/log"
	//"github.com/muaj07/transport/types"
	//"encoding/gob"
	"bytes"
	"os"
	"net"
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
		opts.Logger = log.With(opts.Logger, "Addr", opts.ID)
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


    // Create a new server instance.
    s := &Server{
		TCPTransport: NewTCPTransport(opts.ListenAddr, peerCh),
		peerMap:	make(map[net.Addr]*TCPPeer),
        ServerOpts: opts,
		chain : chain,
        memPool:    NewTxPool(100), //Initial size of the memPool set to 100 txs
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
	case *GetStatusMessage:
		//return s.processsGetStatusMessage(msg.From, tx)
	case *StatusMessage:
		//return s.processStatusMessage(msg.From, tx)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, tx)

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

// func (s *Server) SendGetStatusMessage(tr Transport) error{
// 	var(
// 		getStatusMsg = new(GetStatusMessage)
// 		buf = new(bytes.Buffer)
// 	)

// 	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err !=nil{
// 		return nil
// 	}

// 	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
// 	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err!=nil{
// 		return err
// 	}
// 	return nil
// }

// func (s *Server) processsGetStatusMessage(from net.Addr, data *GetStatusMessage) error{

// 	statusMessage := &StatusMessage{
// 		CurrentHeight: s.chain.Height(),
// 		ID: s.ID,
// 	}
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(statusMessage); err !=nil{
// 		return nil
// 	}
// 	msg := NewMessage(MessageTypeStatus, buf.Bytes())
// 	return s.Transport.SendMessage(from, msg.Bytes())
// }


// func (s *Server) processStatusMessage(from net.Addr, data *StatusMessage) error {
// 	if data.CurrentHeight <= s.chain.Height(){
// 		s.Logger.Log(
// 			"msg", "Cannot Sync Block_Height due to Low or Equal status",
// 			"theHeight", data.CurrentHeight,
// 			"ourHeight", s.chain.Height(),
// 			"From", from,
// 		)
// 		return nil
// 	}
// 	getBlocksMessage := &GetBlocksMessage{
// 		From: s.chain.Height(),
// 		// if To is set to "0", then the max blocks
// 		//will be returned/sent to the requested node
// 		To: 0,
// 		}
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err !=nil{
// 		return err
// 	}
// 	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
// 	return s.Transport.SendMessage(from, msg.Bytes())
// }

func (s *Server) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error{
	panic("Errrrrrr")
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
    //iterate over each transport in the peerMap
    for netAddr, peer := range s.peerMap {
        if err := peer.Send(Payload); err != nil {
            fmt.Printf("Peer send error ==> addr %s [%s]\n",netAddr,err)
        }
    }
    return nil
}




