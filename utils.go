package main

import (
	"crypto/rand"
	"math/big"
)

func generatePrimeBigInt(min, max *big.Int) *big.Int {
	var digit *big.Int
	diff := new(big.Int).Sub(max, min)
	for digit, _ = rand.Int(rand.Reader, diff); !(digit.Add(digit, min)).ProbablyPrime(10); digit, _ = rand.Int(rand.Reader, diff) {
	}
	return digit
}
