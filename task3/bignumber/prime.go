package bignumber

import (
	"math/big"
	"math/rand"
)

// IsPrime returns true if the number is probably prime, false otherwise
func (n *UnsignedBigNumber) IsPrime() bool {
	bigInt := big.NewInt(0)
	bigInt.SetString(n.Hex(), 16)
	return bigInt.ProbablyPrime(5)
}

// FromBigInt converts a big.Int to an UnsignedBigNumber
func FromBigInt(bigInt *big.Int) *UnsignedBigNumber {
	// Convert the big.Int to a byte slice
	bytes := bigInt.Bytes()
	// Reverse the byte slice
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return &UnsignedBigNumber{Values: bytes}
}

// NextPrime returns the next prime number after n
func (n *UnsignedBigNumber) NextPrime() *UnsignedBigNumber {
	num := n.Clone()
	if num.Values[0]&1 == 0 {
		num = num.Add(NewUnsignedBigNumber("1"))
	}

	bigInt := big.NewInt(0)
	bigInt.SetString(num.Hex(), 16)

	for {
		if bigInt.ProbablyPrime(5) {
			return FromBigInt(bigInt)
		}
		bigInt.Add(bigInt, big.NewInt(2))
		if bigInt.ProbablyPrime(5) {
			return FromBigInt(bigInt)
		}
		bigInt.Add(bigInt, big.NewInt(4))
	}
}

// RandPrime generates a random prime number of the specified length
func RandPrime(length uint64) *UnsignedBigNumber {
	for {
		randBytes := make([]byte, length)
		rand.Read(randBytes)
		randNum := &UnsignedBigNumber{Values: randBytes}
		if !randNum.IsPrime() {
			randNum = randNum.NextPrime()
		}
		if randNum.BitLength() == int(length*8) {
			return randNum
		}
	}
}
