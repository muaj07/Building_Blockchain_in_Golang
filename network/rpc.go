package network

import(
	"github.com/muaj07/transport/core"
	"github.com/sirupsen/logrus"
	"encoding/gob"
	"io"
	"bytes"
	"fmt"
    "net"
)

type RPC struct {
	From net.Addr
	Payload io.Reader
}

type RPCHandler interface {
	// Some kind of Decoder
	HandleRPC (rpc RPC) error
}

type RPCProcessor interface{
	ProcessMessage(*DecodeMessage) error
}

//This code declares a new type "RPCDecodeFunc" which is a function that takes  
//an "RPC" parameter and returns a pointer to a "DecodeMessage" and an error

type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

//DefaultRPCDecodeFunc decodes a payload from an RPC into a message struct. 
//The message struct is then parsed based on its header field. If the header 
//is recognized, the message data is decoded and returned along with its 
//originator address. If the header is not recognized, the function returns 
//an error.
func DefaultRPCDecodeFunc(rpc RPC) (*DecodeMessage, error) {
    // Initialize an empty message struct
    msg := Message{}

    // Decode the payload into the message struct using gob package
    if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
        // If the decoding fails, return an error
        return nil, fmt.Errorf("Failed to Decode Message from --> %s Error %s", rpc.From, err)
    }

    // Log the incoming message
    logrus.WithFields(logrus.Fields{
        "Message Type": msg.Header,
        "From":         rpc.From,
    }).Debug("New Incoming Msg")

    // Check the message type
    switch msg.Header {

    // If the message is of type MessageTypeTx (0x1)
    case MessageTypeTx:
        // Decode the transaction data from the "MessageTypeTx" message
        tx := new(core.Transaction) // Transaction is a struct
        if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
            return nil, err
        }
        // Return the decoded transaction data
        return &DecodeMessage{
            From: rpc.From,
            Data: tx,
        }, nil

    // If the message is of type MessageTypeBlock (0x2)
    case MessageTypeBlock:
        // Decode the block data from the "MessageTypeBlock" message
        b := new(core.Block)
        if err := b.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
            return nil, err
        }
        // Return the decoded block data
        return &DecodeMessage{
            From: rpc.From,
            Data: b,
        }, nil

    // If the message is of type MessageTypeGetStatus (0x5)
    case MessageTypeGetStatus:
        // Return a new GetStatusMessage
        return &DecodeMessage{
            From: rpc.From,
            Data: &GetStatusMessage{},
        }, nil

    // If the message is of type MessageTypeStatus (0x4)
    case MessageTypeStatus:
        // Decode the status message from the message data
        statusMessage := new(StatusMessage)
        if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
            return nil, err
        }
        // Return the decoded status message
        return &DecodeMessage{
            From: rpc.From,
            Data: statusMessage,
        }, nil

    // If the message is of type MessageTypeGetBlocks (0x3)
    case MessageTypeGetBlocks:
        getBlocksMessage:= new(GetBlocksMessage)
        if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocksMessage); err != nil {
            return nil, err
        }
        // Return a new GetBlocksMessage
        return &DecodeMessage{
            From: rpc.From,
            Data: &GetBlocksMessage{
                Start: getBlocksMessage.Start,
                End: getBlocksMessage.End,
            },
        }, nil

    // If the message type is not recognized
    default:
        // Return an error
        return nil, fmt.Errorf("INVALID Message Header %x", msg.Header)
    }
}















