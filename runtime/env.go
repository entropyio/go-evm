package runtime

import (
	"github.com/entropyio/go-evm/chain"
	"github.com/entropyio/go-evm/evm"
)

func NewEnv(cfg *Config) *evm.EVM {
	txContext := evm.TxContext{
		Origin:   cfg.Origin,
		GasPrice: cfg.GasPrice,
	}
	blockContext := evm.BlockContext{
		CanTransfer: chain.CanTransfer,
		Transfer:    chain.Transfer,
		GetHash:     cfg.GetHashFn,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		BaseFee:     cfg.BaseFee,
	}

	return evm.NewEVM(blockContext, txContext, cfg.EVMConfig)
}
