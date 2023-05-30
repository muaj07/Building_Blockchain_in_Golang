package network
// import(
// 	//"github.com/muaj07/transport/core"
// 	//"github.com/muaj07/transport/crypto"
// 	"fmt"
// )
type GetStatusMessage struct {

}
type StatusMessage struct{
	ID string
	Version uint32
	CurrentHeight uint32
}

type GetBlocksMessage struct{
	From uint32
	// if To is set to "0", then the max blocks
	//will be returned/sent to the requested node
	To uint32
}