package network

import(
	"github.com/muaj07/transport/core"
	 "github.com/sirupsen/logrus"
	"encoding/gob"
	"io"
	"bytes"
	"fmt"
)

type MessageType byte

const(
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock MessageType = 0x2
    MessageTypeGetBlocks MessageType = 0x3
)

type RPC struct {
	From NetAddr
	Payload io.Reader
}
type RPCHandler interface {
	// Some kind of Decoder
	HandleRPC (rpc RPC) error
}

type Message struct{
	Header MessageType
	Data	[]byte
}

type DecodeMessage struct{
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

type RPCProcessor interface{
	ProcessMessage(*DecodeMessage) error
}

// NewMessage creates a new Message struct with the given MessageType and data.
func NewMessage(t MessageType, data []byte) *Message {
    // Return a pointer to a new Message struct with the given values.
    return &Message {
        Header: t,
        Data: data,
    }
}

// Bytes returns the gob-encoded byte slice of the Message.
func (msg *Message) Bytes() []byte {
    buf := &bytes.Buffer{} // create a new buffer
    enc := gob.NewEncoder(buf) // create a new encoder that writes to the buffer
    enc.Encode(msg) // encode the Message into the buffer
    return buf.Bytes() // return the byte slice of the buffer
}


func DefaultRPCDecodeFunc(rpc RPC) (*DecodeMessage, error) {
	// Initialize an empty message struct
    msg := Message{}
    // Decode the payload into the message struct
    if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
        return nil, fmt.Errorf("Failed to Decode Message from (%s): %s", rpc.From, err)
    }
	logrus.WithFields(logrus.Fields{
		"Message Type": msg.Header,
		"from": rpc.From,
	}).Debug("New incoming Msg")

    switch msg.Header {
    case MessageTypeTx:
        // Decode the transaction data from the "MessageTypeTx" message
        tx := new(core.Transaction)
        if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
            return nil, err
        }
        // Process the transaction
        return &DecodeMessage{
			From: rpc.From,
			Data: tx,
		}, nil

    case MessageTypeBlock:
        // Decode the block data from the "MessageTypeBlock" message
        b := new(core.Block)
        if err := b.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
            return nil, err
        }
        // Process the block
        return &DecodeMessage{
            From: rpc.From,
            Data: b,
            
        }, nil
    default:
        return nil, fmt.Errorf("Invalid Message header %x", msg.Header)
    }

}






