package core
import(
	"bytes"
	"crypto/sha256"
	//"encoding/hex"
	"github.com/muaj07/transport/types"
	//"encoding/gob"
	"time"
)


func NewBlockFromPrevHeader (prevHeader *Header, txx []*Transaction) (*Block, error){
	dataHash, err := CalculateDataHash(txx)
	if err!=nil{
		return nil, err
	}
	header := &Header{
		Version: 1,
		DataHash: dataHash,
		Height: prevHeader.Height+1,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		TimeStamp:  time.Now().UnixNano(),
}
	return 	NewBlock(header, txx)
}


func CalculateDataHash(txx []*Transaction) (types.Hash, error){
	buf := &bytes.Buffer{}
	for _, tx := range txx {
		if err := tx.Encode(NewGobTxEncoder(buf)); err!= nil{
			return types.Hash{}, err
		}
	}
	hash := sha256.Sum256(buf.Bytes())
	return types.Hash(hash), nil
}
