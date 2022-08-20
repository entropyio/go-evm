package evm

import (
	"github.com/entropyio/go-evm/config"
)

func minSwapStack(n int) int {
	return minStack(n, n)
}
func maxSwapStack(n int) int {
	return maxStack(n, n)
}

func minDupStack(n int) int {
	return minStack(n, n+1)
}
func maxDupStack(n int) int {
	return maxStack(n, n+1)
}

func maxStack(pop, push int) int {
	return int(config.StackLimit) + pop - push
}
func minStack(pops, _ int) int {
	return pops
}
