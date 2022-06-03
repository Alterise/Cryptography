package main

import (
	"crypto/rand"
	"math/big"
)

func main() {
	min := new(big.Int)
	max := new(big.Int)
	min.Exp(big.NewInt(2), big.NewInt(129), nil)
	max.Exp(big.NewInt(2), big.NewInt(130), nil)

	P := generatePrimeBigInt(min, max)
	Q := generatePrimeBigInt(min, max)
	//N := new(big.Int).Mul(P, Q)

	mul := new(big.Int).Mul(
		new(big.Int).Mul(new(big.Int).Sub(P, big.NewInt(1)), new(big.Int).Sub(Q, big.NewInt(1))),
		new(big.Int).Mul(new(big.Int).Add(P, big.NewInt(1)), new(big.Int).Add(Q, big.NewInt(1))),
	)

	e := new(big.Int)
	gcd := new(big.Int)
	one := big.NewInt(1)

	for e, _ = rand.Int(rand.Reader, mul); gcd.GCD(nil, nil, mul, e).Cmp(one) != 0; e, _ = rand.Int(rand.Reader, mul) {
	}

	println(e.Text(10))
	//println(mul.Text(10))
	//
	//println(N.Text(10))
	//println(P.Text(10))
	//println(Q.Text(10))
}
