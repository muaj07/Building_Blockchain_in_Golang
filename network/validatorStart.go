package network

import(
	"time"
	//"github.com/muaj07/transport/crypto"
	"github.com/muaj07/transport/core"
)


func (s *Server) validatorLoop(){
	ticker := time.NewTicker(s.BlockTime)
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
	s.Logger.Log(
		"msg", "Validator start Minting Block",
		"blockTime", s.BlockTime,
	)
	// For now, all the available txs in a mempool are
	// included in the block, which will be update in the future
	tx := s.memPool.Pending()
	//"NewBlockFromPrevHeader" is in the helper.go (core)
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
	go s.BroadcastBlock(block)
	return nil
}