package evm

import (
	"github.com/entropyio/go-evm/common"
	"math/big"
)

// calculates the memory size required for a step
func calcMemSize(off, l *big.Int) *big.Int {
	if l.Sign() == 0 {
		return common.Big0
	}

	return new(big.Int).Add(off, l)
}

// getData returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	return common.RightPadBytes(data[start:end], int(size))
}

// getDataBig returns a slice from the data based on the start and size and pads
// up to size with zero's. This function is overflow safe.
func getDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := common.BigMin(start, dlen)
	e := common.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

// bigUint64 returns the integer casted to a uint64 and returns whether it
// overflowed in the process.
func bigUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), v.BitLen() > 64
}

// toWordSize returns the ceiled word size required for memory expansion.
func toWordSize(size uint64) uint64 {
	if size > common.MaxUint64-31 {
		return common.MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

func allZero(b []byte) bool {
	for _, byteArray := range b {
		if byteArray != 0 {
			return false
		}
	}
	return true
}
