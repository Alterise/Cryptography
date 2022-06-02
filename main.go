package main

import (
	"crypto/rand"
	"math/big"
)

func main() {
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(100), nil)
	n, _ := rand.Int(rand.Reader, max)
	for !n.ProbablyPrime(10) {
		max.Exp(big.NewInt(10), big.NewInt(100), nil)
		n, _ = rand.Int(rand.Reader, max)
		n.Add(n, max.Exp(big.NewInt(10), big.NewInt(100), nil))
	}
	print(n.Text(10))
}
