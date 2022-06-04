package main

import (
	"crypto/rand"
	"math/big"
)

var one = big.NewInt(1)
var two = big.NewInt(2)

func generatePrimeBigInt(min, max *big.Int) *big.Int {
	var digit *big.Int
	diff := new(big.Int).Sub(max, min)
	for digit, _ = rand.Int(rand.Reader, diff); !(digit.Add(digit, min)).ProbablyPrime(10); digit, _ = rand.Int(rand.Reader, diff) {
	}
	return digit
}

func calculateD(P *big.Int) *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(P, big.NewInt(2), nil), big.NewInt(4))
}

func lcm(a, b *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(a, b), new(big.Int).GCD(nil, nil, a, b))
}

func funcS(a, b *big.Int) *big.Int {
	return lcm(new(big.Int).Add(a, big.NewInt(1)), new(big.Int).Add(b, big.NewInt(1)))
}

//type Int struct {
//	V big.Int  // Integer value from 0 through M-1
//	M *big.Int // Modulus for finite field arithmetic
//}
//
//func (i *Int) legendre() int {
//	var Pm1, v big.Int
//	Pm1.Sub(i.M, one)
//	v.Div(&Pm1, two)
//	v.Exp(&i.V, &v, i.M)
//	if v.Cmp(&Pm1) == 0 {
//		return -1
//	}
//	return v.Sign()
//}
