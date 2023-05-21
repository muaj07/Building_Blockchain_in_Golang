package network

import(
	"github.com/muaj07/transport/types"
	"github.com/muaj07/transport/core"
	"sort"
)

type TxPool struct{
	transactions map[types.Hash]*core.Transaction
}

type TxMapSorter struct {
	txx []*core.Transaction
}

// NewTxMapSorter returns a pointer to a new TxMapSorter instance
// created from the given map of transactions
func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter{
	txx := make([]*core.Transaction, len(txMap))

	i:= 0
	for _,tx := range txMap{
		txx[i]=tx
		i++
	}
	s:=&TxMapSorter{txx}
	sort.Sort(s)
	return s
}

func (s *TxMapSorter) Len() int{
	return len(s.txx)
}
func (s *TxMapSorter) Swap(i,j int) {
	s.txx[i], s.txx[j] = s.txx[j],s.txx[i]
}
func (s *TxMapSorter) Less(i,j int) bool{
	return s.txx[i].FirstSeen() < s.txx[j].FirstSeen()
}

func NewTxPool() *TxPool{
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Transactions returns a slice of pointers to core.Transaction objects.
// It sorts the transactions in the TxPool using NewTxMapSorter before returning
// them.
func(t *TxPool) Transactions() []*core.Transaction{
	txs := NewTxMapSorter(t.transactions)
	return txs.txx
}


// AddTx adds a transaction to the transaction pool
func (t *TxPool) AddTx(tx *core.Transaction) error {
    // Get the hash of the transaction
    hash := tx.Hash(core.TxHasher{})
	t.transactions[hash] = tx
    // Return nil to indicate success
    return nil
}

// Has checks if a transaction with a given hash exists in the pool.
// Returns true if the transaction exists, false otherwise.
func (t *TxPool) Has(hash types.Hash) bool{
	_,ok := t.transactions[hash]
	return ok
}

func (t *TxPool) Len() int{
	return len(t.transactions)	
}

// Flush clears all transactions from the transaction pool.
func (t *TxPool) Flush() {
    // Create a new empty map to replace the current map of transactions.
    t.transactions = make(map[types.Hash]*core.Transaction)
}
