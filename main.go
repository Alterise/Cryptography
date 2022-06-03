package main

import "math/big"

func main() {
	min := new(big.Int)
	max := new(big.Int)
	min.Exp(big.NewInt(10), big.NewInt(100), nil)
	max.Exp(big.NewInt(10), big.NewInt(100), nil)
	max.Mul(max, big.NewInt(2))
	digit := generatePrimeBigInt(min, max)
	print(digit.Text(10))
}
