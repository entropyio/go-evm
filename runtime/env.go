package runtime

import (
	"github.com/entropyio/go-evm/evm"
)

func NewEnv(cfg *Config) *evm.EVM {
	context := evm.Context{
		Origin:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		GasPrice:    cfg.GasPrice,
	}

	return evm.NewEVM(context, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
}
