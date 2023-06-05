package network

import(
	"bytes"
	"time"
    "net"
	"fmt"
	"encoding/gob"
	//"github.com/muaj07/transport/transport"
	"github.com/muaj07/transport/core"
)

// ProcessMessage processes a Decoded Message received by the server
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
		return s.processsGetStatusMessage(msg.From, tx)
	case *StatusMessage:
		return s.processStatusMessage(msg.From, tx)
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
	go s.BroadcastBlock(b)
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
	go s.BroadcastTx(tx)
	s.memPool.Add(tx)
	return nil
}

func (s *Server) SendGetStatusMessage(peer *TCPPeer) error{
	var(
		getStatusMsg = new(GetStatusMessage)
		buf = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err !=nil{
		return nil
	}
	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	s.Logger.Log(
		"msg", "Sent GetStatusMessage",
        "peer", peer,
	)
	return peer.Send(msg.Bytes())
}

func (s *Server) processsGetStatusMessage(from net.Addr, data *GetStatusMessage) error{
	s.Logger.Log(
		"msg", "Received GetStatusMessage",
        "from", from,
	)

	statusMessage := &StatusMessage{
		CurrentHeight: s.chain.Height(),
		ID: s.ID,
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err !=nil{
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	peer,ok := s.peerMap[from]
	if !ok{
		return fmt.Errorf("Peer %s unknown", peer.conn.RemoteAddr())
		}
	msg := NewMessage(MessageTypeStatus, buf.Bytes())
	return peer.Send(msg.Bytes())
}


func (s *Server) processStatusMessage(from net.Addr, data *StatusMessage) error {
	s.Logger.Log(
		"msg", "Received Status Message",
        "from", from,
	)
	if data.CurrentHeight <= s.chain.Height(){
		s.Logger.Log(
			"msg", "Cannot Sync Block_Height due to Low or Equal status",
			"theHeight", data.CurrentHeight,
			"ourHeight", s.chain.Height(),
			"From", from,
		)
		return nil
	}
	getBlocksMessage := &GetBlocksMessage{
		Start: s.chain.Height(),
		// if To is set to "0", then the max blocks
		//will be returned/sent to the requested node
		End: 0,
		}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err !=nil{
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	peer,ok := s.peerMap[from]
	if !ok{
		return fmt.Errorf("Peer %s unknown", peer.conn.RemoteAddr())
		}
	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
	return peer.Send(msg.Bytes())
}



func (s *Server) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error{
	panic("Err")
	return nil
}






