package network
import(
	//"github.com/muaj07/transport/core"
	//"github.com/muaj07/transport/crypto"
	"bytes"
	"net"
	"encoding/gob"
)

type MessageType byte

const(
	MessageTypeTx MessageType = 0x1 //Tx msg 
	MessageTypeBlock MessageType = 0x2 //Block msg
    MessageTypeGetBlocks MessageType = 0x3
    MessageTypeStatus MessageType = 0x4
    MessageTypeGetStatus MessageType = 0x5
)
type Message struct{
	Header MessageType
	Data	[]byte
}

type DecodeMessage struct{
	From net.Addr
	Data any
}

type GetStatusMessage struct {

}

type StatusMessage struct{
	ID string
	Version uint32
	CurrentHeight uint32
}

type GetBlocksMessage struct{
	Start uint32
	// if To is set to "0", then the max blocks
	//will be returned/sent to the requested node
	End uint32
}



// NewMessage creates a new Message struct with the given MessageType and data.
func NewMessage(t MessageType, data []byte) *Message {
    // Return a pointer to a new Message struct with the given values.
    return &Message {
        Header: t, //Message type 
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