package network
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/muaj07/transport/core"
	"strconv"
	"math/rand"
	//"math/big"
)

func TestTxPool(t *testing.T) {
	
	p := NewTxPool()
	assert.Equal(t, p.Len(),0)
}

func TestTxPool_Add(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("Add Trx"))
	assert.Nil(t,p.AddTx(tx))
	assert.Equal(t, p.Len(), 1)
	_ = core.NewTransaction([]byte("Add Trx"))
	assert.Equal(t, p.Len(),1)

	p.Flush()
	assert.Equal(t, p.Len(),0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txlen :=100

	for i:=0; i<txlen; i++{
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i),10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t,p.AddTx(tx))
	}
	assert.Equal(t, txlen, p.Len())

	txx := p.Transactions()
	for i:=0; i<len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
	
}