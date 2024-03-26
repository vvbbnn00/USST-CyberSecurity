package bignumber

import (
	"math/big"
	"math/rand"
)

// IsPrime returns true if the number is probably prime, false otherwise
func (n *BigNumber) IsPrime() bool {
	bigInt := big.NewInt(0)
	bigInt.SetString(n.Hex(), 16)
	return bigInt.ProbablyPrime(5)
}

// NextPrime returns the next prime number after n
func (n *BigNumber) NextPrime() *BigNumber {
	num := n.Clone()
	if num.AbsVal.Values[0]&1 == 0 {
		num = num.Add(NewBigNumber("1"))
	}

	bigInt := big.NewInt(0)
	bigInt.SetString(num.Hex(), 16)

	for {
		if bigInt.ProbablyPrime(5) {
			return FromBigIntSigned(bigInt)
		}
		bigInt.Add(bigInt, big.NewInt(2))
		if bigInt.ProbablyPrime(5) {
			return FromBigIntSigned(bigInt)
		}
		bigInt.Add(bigInt, big.NewInt(4))
	}
}

// RandPrime generates a random prime number of the specified length
func RandPrime(length uint64) *BigNumber {
	for {
		randBytes := make([]byte, length)
		rand.Read(randBytes)
		randNum := &BigNumber{Sign: 1, AbsVal: &UnsignedBigNumber{Values: randBytes}}
		if !randNum.IsPrime() {
			randNum = randNum.NextPrime()
		}
		if randNum.BitLength() == int(length*8) {
			return randNum
		}
	}
}
