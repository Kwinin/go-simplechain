package cross

import (
	"math/big"
	"time"

	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/log"
)

type CtxStoreConfig struct {
	ChainId      *big.Int
	Anchors      []common.Address
	IsAnchor     bool
	Rejournal    time.Duration // Time interval to regenerate the local transaction journal
	ValueLimit   *big.Int      // Minimum value to enforce for acceptance into the pool
	AccountSlots uint64        // Number of executable transaction slots guaranteed per account
	GlobalSlots  uint64        // Maximum number of executable transaction slots for all accounts
	AccountQueue uint64        // Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64        // Maximum number of non-executable transaction slots for all accounts
}

var DefaultCtxStoreConfig = CtxStoreConfig{
	Anchors:      []common.Address{},
	Rejournal:    time.Minute * 10,
	ValueLimit:   big.NewInt(1e18),
	AccountSlots: 5,
	GlobalSlots:  4096,
	AccountQueue: 5,
	GlobalQueue:  10,
}

func (config *CtxStoreConfig) Sanitize() CtxStoreConfig {
	conf := *config
	if conf.Rejournal < time.Second {
		log.Warn("Sanitizing invalid ctxpool journal time", "provided", conf.Rejournal, "updated", time.Second)
		conf.Rejournal = time.Second
	}
	if conf.ValueLimit == nil || conf.ValueLimit.Cmp(big.NewInt(1e18)) < 0 {
		log.Warn("Sanitizing invalid ctxpool price limit", "provided", conf.ValueLimit, "updated", DefaultCtxStoreConfig.ValueLimit)
		conf.ValueLimit = DefaultCtxStoreConfig.ValueLimit
	}
	return conf
}
