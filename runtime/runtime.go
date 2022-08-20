package runtime

import (
	"github.com/entropyio/go-evm/common"
	"github.com/entropyio/go-evm/common/crypto"
	"github.com/entropyio/go-evm/config"
	"github.com/entropyio/go-evm/evm"
	"github.com/entropyio/go-evm/logger"
	"math"
	"math/big"
	"time"
)

var log = logger.NewLogger("[runtime]")

// Config is a basic type specifying certain configuration flags for running
// the EVM.
type Config struct {
	ChainConfig *config.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   evm.EVMConfig
	BaseFee     *big.Int

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults(cfg *Config) {
	//if cfg.ChainConfig == nil {
	//	cfg.ChainConfig = &config.ChainConfig{
	//		ChainID:        big.NewInt(1),
	//		HomesteadBlock: new(big.Int),
	//		EIP150Block:    new(big.Int),
	//		EIP155Block:    new(big.Int),
	//		EIP158Block:    new(big.Int),
	//	}
	//}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
	if cfg.BaseFee == nil {
		cfg.BaseFee = big.NewInt(config.InitialBaseFee)
	}
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Execute sets up an in-memory, temporary, environment for the execution of
// the given code. It makes sure that it's restored to its original state afterwards.
func Execute(code, input []byte, cfg *Config) ([]byte, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	//if cfg.State == nil {
	//	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	//}
	var (
		address = common.BytesToAddress([]byte("contract"))
		vmenv   = NewEnv(cfg)
		sender  = evm.AccountRef(cfg.Origin)
	)
	//if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil); rules.IsBerlin {
	//	cfg.State.PrepareAccessList(cfg.Origin, &address, vm.ActivePrecompiles(rules), nil)
	//}
	//cfg.State.CreateAccount(address)
	// set the receiver's (the executing contract) code for execution.
	//cfg.State.SetCode(address, code)
	log.Debugf("execute address:%x, code:%+v, input:%+v", address, code, input)

	// Call the code with the given configuration.
	ret, _, err := vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, err
}

// Create executes the code using the EVM create method
func Create(input []byte, cfg *Config) ([]byte, common.Address, uint64, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	//if cfg.State == nil {
	//	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	//}
	var (
		vmenv  = NewEnv(cfg)
		sender = evm.AccountRef(cfg.Origin)
	)
	//if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil); rules.IsBerlin {
	//	cfg.State.PrepareAccessList(cfg.Origin, nil, vm.ActivePrecompiles(rules), nil)
	//}
	// Call the code with the given configuration.
	code, address, leftOverGas, err := vmenv.Create(
		sender,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return code, address, leftOverGas, err
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
//
// Call, unlike Execute, requires a config and also requires the State field to
// be set.
func Call(address common.Address, input []byte, cfg *Config) ([]byte, uint64, error) {
	setDefaults(cfg)

	vmenv := NewEnv(cfg)
	sender := evm.AccountRef(cfg.Origin) // TODO: Call/Excuter sender not the same
	//sender := cfg.State.GetOrNewStateObject(cfg.Origin)

	//statedb := cfg.State

	//if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil); rules.IsBerlin {
	//	statedb.PrepareAccessList(cfg.Origin, &address, vm.ActivePrecompiles(rules), nil)
	//}
	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return ret, leftOverGas, err
}
