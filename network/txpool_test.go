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
	
	p := NewTxPool(100)
	assert.Equal(t, len(p.Pending()),0)
}

func TestTxPool_Add(t *testing.T) {
	p := NewTxPool(100)
	tx := core.NewTransaction([]byte("Add Trx"))
	p.Add(tx)
	//assert.Nil(t,p.Add(tx))
	assert.Equal(t, len(p.Pending()), 1)
	_ = core.NewTransaction([]byte("Add Trx"))
	assert.Equal(t, len(p.Pending()),1)

	p.ClearPending()
	assert.Equal(t, len(p.Pending()),0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool(100)
	txlen :=10

	for i:=0; i<txlen; i++{
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i),10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		p.Add(tx)
	}
	assert.Equal(t, txlen, len(p.Pending()))

	// tx := p.PendingCount()
	// for i:=0; i<tx-1; i++ {
	// 	assert.True(t, tx[i].FirstSeen() < tx[i+1].FirstSeen())
	// }
	
}