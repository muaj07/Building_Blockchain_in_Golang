package network

import(
	"github.com/muaj07/transport/types"
	"github.com/muaj07/transport/core"
	"sync"
)

type TxSortedMap struct{
	lock sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txx *types.List[*core.Transaction]
}

type TxPool struct{
	all 	*TxSortedMap
	pending *TxSortedMap
	//The max length of the total pool
	//when the pool is full, remove the oldest trx to
	//create space for the new trxs
	maxLength int
}



func NewTxPool(maxLength int) *TxPool{
	return &TxPool{
		all: NewTxSortedMap(),
		pending: NewTxSortedMap(),
		maxLength: maxLength,
	}
}

// AddTx adds a transaction to the transaction pool
func (t *TxPool) Add(tx *core.Transaction) {
    // if the mempool if full 
	// prune the oldest transaction sitting in the txPool

	if t.all.Count() == t.maxLength {
		oldest := t.all.First()
		t.all.Remove(oldest.Hash(core.TxHasher{}))
	}
	// first check if the tx is not already in the mempool
	if !t.all.Contains(tx.Hash(core.TxHasher{})) {
		t.all.Add(tx)
		t.pending.Add(tx)
	}
}


func (t *TxPool) Contains (hash types.Hash) bool{
	return t.all.Contains(hash) // call the Contains method for NewTxSortedMap and return bool
}

func (t *TxPool) Pending() []*core.Transaction{
	return t.pending.txx.Data //return the pending txs in the mempool (i.e. TxPool)
}

func (t *TxPool) ClearPending() {
	t.pending.Clear() //return the number of pending txs
}

func (t *TxPool) PendingCount() int {
	return t.pending.Count() //return the number of pending txs
}



// NewTxSortedMap returns a pointer to a new TxSortedMap instance
// created from the given map of transactions
func NewTxSortedMap() *TxSortedMap{
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txx: types.NewList[*core.Transaction](),
	}
}

func( s *TxSortedMap) First() *core.Transaction {
	s.lock.RLock()
	defer s.lock.RUnlock()
	first := s.txx.Get(0)
	return s.lookup[first.Hash(core.TxHasher{})]
}

func (s *TxSortedMap) Get(h types.Hash) *core.Transaction{
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.lookup[h]
}



func (s *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.lookup[hash]; !ok{
		s.lookup[hash]=tx
		s.txx.Insert(tx)
	}

}

func (s *TxSortedMap) Remove(h types.Hash) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.txx.Remove(s.lookup[h])
	delete(s.lookup, h)
}

func (s *TxSortedMap) Count() int{
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.lookup)
}



func(s *TxSortedMap) Contains(h types.Hash) bool{
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.lookup[h]
	return ok
}



func (s *TxSortedMap) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lookup = make(map[types.Hash]*core.Transaction)
	s.txx.Clear()
}

