package network

import(
	"github.com/muaj07/transport/core"
	//"github.com/sirupsen/logrus"
	"encoding/gob"
	"io"
	"bytes"
	"fmt"
)

type MessageType byte

const(
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock
)

type RPC struct {
	From NetAddr
	Payload io.Reader
}

type Message struct{
	Header MessageType
	Data	[]byte
}

type RPCHandler interface {
	// Some kind of Decoder
	HandleRPC (rpc RPC) error
}

type RPCProcessor interface{
	ProcessTransaction(NetAddr, *core.Transaction) error
}

type DefaultRPCHandler struct{
	p RPCProcessor
}


func NewDefaultRPCHandler (p RPCProcessor) *DefaultRPCHandler{
	return &DefaultRPCHandler {
		p: p,
	}
}

// HandleRPC handles the incoming RPC request
// and returns an error if any.
func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
    // Initialize an empty message struct
    msg := Message{}
    // Decode the payload into the message struct
    if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
        return fmt.Errorf("Failed to Decode Message from (%s): %s", rpc.From, err)
    }

    switch msg.Header {
    case MessageTypeTx:
        // Decode the transaction data from the message
        tx := new(core.Transaction)
        if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
            return err
        }
        // Process the transaction
        return h.p.ProcessTransaction(rpc.From, tx)
    default:
        return fmt.Errorf("Invalid Message header %x", msg.Header)
    }
    return nil
}



func (p *DefaultRPCHandler) ProcessTransaction(from NetAddr, tx *core.Transaction) error{
	return nil
}

func NewMessage (t MessageType, data []byte) *Message{
	return &Message{
		Header: t,
		Data: data,
	}
}
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(msg)
	return buf.Bytes()
}