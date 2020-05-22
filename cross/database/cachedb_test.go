package db

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/simplechain-org/go-simplechain/common"
	cc "github.com/simplechain-org/go-simplechain/cross/core"
	"github.com/simplechain-org/go-simplechain/ethdb/memorydb"
)

func TestCtxDb(t *testing.T) {
	var (
		db    = memorydb.New()
		ctxDb = NewCacheDb(big.NewInt(1), db, 10)
	)
	var i int64
	for i = 0; i < 1024; i++ {
		ctx := cc.NewCrossTransactionWithSignatures(cc.NewCrossTransaction(big.NewInt(17),
			big.NewInt(rand.Int63n(110)),
			big.NewInt(1),
			common.BigToHash(big.NewInt(i)),
			common.Hash{},
			common.Hash{},
			common.Address{},
			nil), 0)
		err := ctxDb.Write(ctx)
		if err != nil {
			t.Error(err)
		}
		if !ctxDb.Has(ctx.ID()) {
			t.Errorf("write err,id:%s", ctx.ID().String())
		}
	}

	if len(ctxDb.Query(0, 0)) != 1024 {
		t.Errorf("write count err,len:%d", len(ctxDb.Query(0, 0)))
	}

	cws, err := ctxDb.Read(common.BigToHash(big.NewInt(1000)))
	if err != nil {
		t.Error(err)
	}
	if err := ctxDb.Delete(cws.ID()); err != nil {
		t.Error(err)
	}

	if ctxDb.Has(cws.ID()) {
		t.Errorf("Delete err,id:%s", cws.ID().String())
	}

	if len(ctxDb.Query(0, 0)) != 1023 {
		t.Errorf("write count err,len:%d", len(ctxDb.Query(0, 0)))
	}
}
