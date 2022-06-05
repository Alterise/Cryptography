package main

import (
	"crypto/rand"
	"math/big"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)

type LucKey struct {
	keyValue []byte
	mulValue []byte
}

func generatePrimeBigInt(min, max *big.Int) *big.Int {
	var digit *big.Int
	diff := new(big.Int).Sub(max, min)
	for digit, _ = rand.Int(rand.Reader, diff); !(digit.Add(digit, min)).ProbablyPrime(100); digit, _ = rand.Int(rand.Reader, diff) {
	}
	return digit
}

func calculateD(P *big.Int) *big.Int {
	return new(big.Int).Sub(new(big.Int).Exp(P, two, nil), big.NewInt(4))
}

func lcm(a, b *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(a, b), new(big.Int).GCD(nil, nil, a, b))
}

func funcS(p, q, Dp, Dq *big.Int) *big.Int {
	return lcm(new(big.Int).Sub(p, Dp), new(big.Int).Sub(q, Dq))
}

func legendre(a, p *big.Int) *big.Int {
	if new(big.Int).Mod(a, p) == zero {
		return big.NewInt(0)
	}

	if isQuadraticResidue(a, p) {
		return big.NewInt(1)
	} else {
		return big.NewInt(-1)
	}

}

func isQuadraticResidue(a, p *big.Int) bool {
	return new(big.Int).Exp(a, new(big.Int).Div(new(big.Int).Sub(p, one), two), p).Cmp(one) == 0
}

func calculateLucasSlow(n, p, q *big.Int) *big.Int {
	if n.Cmp(zero) == 0 {
		return big.NewInt(2)
	}
	if n.Cmp(one) == 0 {
		return new(big.Int).Set(p)
	}
	if new(big.Int).Mod(n, two).Cmp(zero) == 0 {
		currN := new(big.Int).Div(n, two)
		return new(big.Int).Sub(new(big.Int).Exp(calculateLucasSlow(currN, p, q), two, nil), new(big.Int).Mul(two, new(big.Int).Exp(q, currN, nil)))
	} else {
		currM := new(big.Int).Div(n, two)
		currN := new(big.Int).Add(currM, one)
		return new(big.Int).Sub(new(big.Int).Mul(calculateLucasSlow(currN, p, q), calculateLucasSlow(currM, p, q)), new(big.Int).Mul(new(big.Int).Exp(q, currN, nil), p))
	}
}

func calculateCount(n *big.Int) []*big.Int {
	count := []*big.Int{big.NewInt(-1)}
	id := new(big.Int).Set(n)

	for id.Cmp(one) == 1 {
		if new(big.Int).Mod(id, two).Cmp(one) == 0 {
			count = append(count, big.NewInt(1))
		} else {
			count = append(count, big.NewInt(0))
		}
		id.Div(id, two)
	}
	return count
}

func calculateLucas(n, P, m []byte) *big.Int {
	rems := calculateCount(new(big.Int).SetBytes(n))
	prev := big.NewInt(2)
	curr := new(big.Int).SetBytes(P)
	PInternal := new(big.Int).SetBytes(P)
	mInternal := new(big.Int).SetBytes(m)

	for i := len(rems) - 1; i >= 0; i-- {
		currId := rems[i]
		if currId.Cmp(zero) == 0 {
			Ve := new(big.Int).Sub(new(big.Int).Mul(curr, curr), two)
			Vo := new(big.Int).Sub(new(big.Int).Mul(curr, prev), PInternal)
			prev = new(big.Int).Mod(Vo, mInternal)
			curr = new(big.Int).Mod(Ve, mInternal)
		} else if currId.Cmp(one) == 0 {
			Vo := new(big.Int).Sub(new(big.Int).Sub(new(big.Int).Mul(new(big.Int).Mul(curr, curr), PInternal), new(big.Int).Mul(curr, prev)), PInternal)
			prev = new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(curr, curr), two), mInternal)
			curr = new(big.Int).Mod(Vo, mInternal)
		} else {
			break
		}
	}
	return curr
}

func generateKeys(data []byte) (LucKey, LucKey) {
	dataInternal := new(big.Int).SetBytes(data)
	min := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(len(data)*4)), nil)
	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(len(data)*4)+1), nil)

	p := generatePrimeBigInt(min, max)
	q := generatePrimeBigInt(min, max)

	N := new(big.Int).Mul(p, q)

	mul := new(big.Int).Mul(
		new(big.Int).Mul(new(big.Int).Sub(p, one), new(big.Int).Sub(q, one)),
		new(big.Int).Mul(new(big.Int).Add(p, one), new(big.Int).Add(q, one)),
	)

	e := new(big.Int)
	gcd := new(big.Int)

	for e, _ = rand.Int(rand.Reader, N); gcd.GCD(nil, nil, mul, e).Cmp(one) != 0; e, _ = rand.Int(rand.Reader, N) {
	}

	D := calculateD(dataInternal)
	Dp := legendre(D, p)
	Dq := legendre(D, q)
	S := funcS(p, q, Dp, Dq)

	d := new(big.Int).ModInverse(e, S)
	return LucKey{e.Bytes(), N.Bytes()}, LucKey{d.Bytes(), N.Bytes()}
}

func encDec(key LucKey, message []byte) []byte {
	return calculateLucas(key.keyValue, message, key.mulValue).Bytes()
}
