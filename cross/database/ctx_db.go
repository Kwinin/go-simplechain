package db

import (
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/core/rawdb"
	cc "github.com/simplechain-org/go-simplechain/cross/core"
	"github.com/simplechain-org/go-simplechain/ethdb"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type ErrCtxDbFailure struct {
	msg string
	err error
}

func (e ErrCtxDbFailure) Error() string {
	return fmt.Sprintf("DB ctx handle failed:%s %s", e.msg, e.err)
}

type ServiceContext interface {
	ResolvePath(string) string
	OpenDatabase(string, int, int, string) (ethdb.Database, error)
}

type CtxDB interface {
	io.Closer
	ChainID() *big.Int
	//Size() int
	Count(filter ...q.Matcher) int
	Height() uint64
	Load() error
	Repair() error
	Write(ctx *cc.CrossTransactionWithSignatures) error
	Read(ctxId common.Hash) (*cc.CrossTransactionWithSignatures, error)
	One(field FieldName, key interface{}) *cc.CrossTransactionWithSignatures
	//ReadAll(ctxId common.Hash) (common.Hash, error)
	//Delete(ctxId common.Hash, atNumber uint64) error
	Update(id common.Hash, updater func(ctx *CrossTransactionIndexed)) error
	Has(id common.Hash) bool
	Query(pageSize int, startPage int, orderBy FieldName, filter ...q.Matcher) []*cc.CrossTransactionWithSignatures
	RangeByNumber(begin, end uint64, pageSize int) []*cc.CrossTransactionWithSignatures

	//QueryByPK(pageSize int, startPage int, filter ...q.Matcher) []*cc.CrossTransactionWithSignatures
	//QueryByPrice(pageSize int, startPage int, filter ...q.Matcher) []*cc.CrossTransactionWithSignatures
	//Range(pageSize int, startCtxID *common.Hash, endCtxID *common.Hash) []*cc.CrossTransactionWithSignatures
}

func OpenStormDB(ctx ServiceContext, name string) (*storm.DB, error) {
	if ctx == nil || len(ctx.ResolvePath(name)) == 0 {
		return storm.Open(os.TempDir() + name)
	}
	return storm.Open(ctx.ResolvePath(name))
}

func OpenEtherDB(ctx ServiceContext, name string) (ethdb.Database, error) {
	if ctx == nil {
		return rawdb.NewMemoryDatabase(), nil
	}
	//TODO: cache & handles
	db, err := ctx.OpenDatabase(name, 16, 256, "")
	if err != nil {
		return nil, ErrCtxDbFailure{msg: "OpenEtherDB fail", err: err}
	}
	return db, nil
}
