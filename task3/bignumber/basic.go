package bignumber

import (
	"fmt"
	"math"
	"math/big"
)

// UnsignedBigNumber represents an unsigned big number
type UnsignedBigNumber struct {
	Values []byte // 0x12345678 -> [0x78, 0x56, 0x34, 0x12]
}

// NewUnsignedBigNumber creates a new UnsignedBigNumber from a hex string
func NewUnsignedBigNumber(hexStr string) *UnsignedBigNumber {
	n := &UnsignedBigNumber{
		Values: make([]byte, (len(hexStr)+1)/2),
	}
	for i, j := len(hexStr)-1, 0; i >= 0; i, j = i-2, j+1 {
		if i == 0 {
			n.Values[j] = byteFromHexChar(hexStr[i])
		} else {
			n.Values[j] = byteFromHexChar(hexStr[i-1])<<4 | byteFromHexChar(hexStr[i])
		}
	}
	return n.trim()
}

// byteFromHexChar converts a hex character to a byte
func byteFromHexChar(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return 10 + c - 'a'
	case 'A' <= c && c <= 'F':
		return 10 + c - 'A'
	default:
		panic("invalid hex character")
	}
}

// Hex converts the UnsignedBigNumber to a hex string
func (n *UnsignedBigNumber) Hex() string {
	var result string
	for i := len(n.Values) - 1; i >= 0; i-- {
		result += fmt.Sprintf("%02x", n.Values[i])
	}
	if result == "" {
		result = "00"
	}
	return result
}

// Add adds two UnsignedBigNumbers and returns the result
func (n *UnsignedBigNumber) Add(other *UnsignedBigNumber) *UnsignedBigNumber {
	maxLen := int(math.Max(float64(len(n.Values)), float64(len(other.Values))))
	result := &UnsignedBigNumber{Values: make([]byte, maxLen+1)}
	var carry uint64
	for i := 0; i < maxLen; i++ {
		var a, b uint64
		if i < len(n.Values) {
			a = uint64(n.Values[i])
		}
		if i < len(other.Values) {
			b = uint64(other.Values[i])
		}
		sum := a + b + carry
		result.Values[i] = uint8(sum & 0xFF)
		carry = sum >> 8
	}
	if carry > 0 {
		result.Values[maxLen] = uint8(carry & 0xFF)
	}
	return result
}

// Subtract subtracts one UnsignedBigNumber from another and returns the result
func (n *UnsignedBigNumber) Subtract(other *UnsignedBigNumber) *UnsignedBigNumber {
	if n.LessThan(other) {
		panic("negative result not allowed for UnsignedBigNumber")
	}
	result := &UnsignedBigNumber{Values: make([]byte, len(n.Values))}
	var borrow uint64
	for i := 0; i < len(n.Values); i++ {
		var a, b uint64
		a = uint64(n.Values[i])
		if i < len(other.Values) {
			b = uint64(other.Values[i])
		}
		if a < b+borrow {
			result.Values[i] = byte(a + 0x100 - b - borrow)
			borrow = 1
		} else {
			result.Values[i] = byte(a - b - borrow)
			borrow = 0
		}
	}
	return result.trim()
}

// SimpleMultiply multiplies two UnsignedBigNumbers and returns the result, use multi-threading to speed up
func (n *UnsignedBigNumber) SimpleMultiply(other *UnsignedBigNumber) *UnsignedBigNumber {
	if n.IsZero() || other.IsZero() {
		return &UnsignedBigNumber{Values: []byte{0}}
	}

	result := &UnsignedBigNumber{Values: make([]byte, len(n.Values)+len(other.Values))}
	for i := 0; i < len(n.Values); i++ {
		for j := 0; j < len(other.Values); j++ {
			mul := uint64(n.Values[i]) * uint64(other.Values[j])
			// add the product to the result
			k := i + j
			for mul > 0 {
				mul += uint64(result.Values[k])
				result.Values[k] = byte(mul & 0xFF)
				mul >>= 8
				k++
			}
		}
	}
	return result.trim()
}

// BitLength returns the number of bin in the UnsignedBigNumber
func (n *UnsignedBigNumber) BitLength() int {
	if n.IsZero() {
		return 0
	}
	result := (len(n.Values) - 1) * 8
	v := n.Values[len(n.Values)-1]
	for v > 0 {
		result++
		v >>= 1
	}
	return result
}

func (n *UnsignedBigNumber) Clone() *UnsignedBigNumber {
	return &UnsignedBigNumber{Values: append([]byte{}, n.Values...)}
}

// SimpleDivide divides one UnsignedBigNumber by another and returns the quotient and remainder
func (n *UnsignedBigNumber) SimpleDivide(other *UnsignedBigNumber) (*UnsignedBigNumber, *UnsignedBigNumber) {
	if other.IsZero() {
		panic("division by zero")
	}
	if n.LessThan(other) {
		return &UnsignedBigNumber{Values: []byte{0}}, n.Clone()
	}
	quotient := &UnsignedBigNumber{Values: make([]byte, len(n.Values)-len(other.Values)+1)}
	remainder := n.Clone()
	for i := len(n.Values) - len(other.Values); i >= 0; i-- {
		shift := 0
		for !remainder.LessThan(other.ShiftLeft(uint(i * 8))) {
			shift++
			remainder = remainder.Subtract(other.ShiftLeft(uint(i * 8)))
		}
		quotient.Values[i] = byte(shift)
	}
	return quotient.trim(), remainder
}

func (n *UnsignedBigNumber) Divide(other *UnsignedBigNumber) (*UnsignedBigNumber, *UnsignedBigNumber) {
	//start := time.Now()
	quotient, remainder := n.SimpleDivide(other)
	//fmt.Println("Divide:", time.Since(start))
	return quotient, remainder
}

// IsZero checks if the UnsignedBigNumber is zero
func (n *UnsignedBigNumber) IsZero() bool {
	for _, v := range n.Values {
		if v != 0 {
			return false
		}
	}
	return true
}

// GreaterThanOrEqual checks if one UnsignedBigNumber is greater than or equal to another
func (n *UnsignedBigNumber) GreaterThanOrEqual(other *UnsignedBigNumber) bool {
	if len(n.Values) > len(other.Values) {
		return true
	} else if len(n.Values) < len(other.Values) {
		return false
	}
	for i := len(n.Values) - 1; i >= 0; i-- {
		if n.Values[i] > other.Values[i] {
			return true
		} else if n.Values[i] < other.Values[i] {
			return false
		}
	}
	return true
}

// ShiftLeft shifts the UnsignedBigNumber left by a number of bits
func (n *UnsignedBigNumber) ShiftLeft(bits uint) *UnsignedBigNumber {
	result := n.Clone()

	remainBits := bits
	for {
		if remainBits == 0 {
			break
		}

		shift := remainBits % 8
		if shift == 0 {
			shift = 8
		}

		var carry byte
		for i := 0; i < len(result.Values); i++ {
			newCarry := result.Values[i] >> (8 - shift)
			result.Values[i] = (result.Values[i] << shift) | carry
			carry = newCarry
		}

		if carry > 0 {
			result.Values = append(result.Values, carry)
		}

		remainBits -= shift
	}

	return result
}

// ShiftRight shifts the UnsignedBigNumber right by a number of bits
func (n *UnsignedBigNumber) ShiftRight(bits uint) *UnsignedBigNumber {
	result := n.Clone()

	remainBits := bits
	for {
		if remainBits == 0 {
			break
		}

		shift := remainBits % 8
		if shift == 0 {
			shift = 8
		}

		var carry byte
		for i := len(result.Values) - 1; i >= 0; i-- {
			newCarry := result.Values[i] << (8 - shift)
			result.Values[i] = (result.Values[i] >> shift) | carry
			carry = newCarry
		}

		if result.Values[len(result.Values)-1] == 0 {
			result.Values = result.Values[:len(result.Values)-1]
		}

		remainBits -= shift
	}

	return result
}

// trim removes leading zeros from the UnsignedBigNumber
func (n *UnsignedBigNumber) trim() *UnsignedBigNumber {
	result := n.Clone()
	i := len(result.Values) - 1
	for i >= 0 && result.Values[i] == 0 {
		i--
	}
	result.Values = result.Values[:i+1]
	return result
}

// Equals checks if two UnsignedBigNumbers are equal
func (n *UnsignedBigNumber) Equals(other *UnsignedBigNumber) bool {
	if len(n.Values) != len(other.Values) {
		return false
	}
	for i := range n.Values {
		if n.Values[i] != other.Values[i] {
			return false
		}
	}
	return true
}

// LessThan checks if one UnsignedBigNumber is less than another
func (n *UnsignedBigNumber) LessThan(other *UnsignedBigNumber) bool {
	if len(n.Values) < len(other.Values) {
		return true
	} else if len(n.Values) > len(other.Values) {
		return false
	}
	for i := len(n.Values) - 1; i >= 0; i-- {
		if n.Values[i] < other.Values[i] {
			return true
		} else if n.Values[i] > other.Values[i] {
			return false
		}
	}
	return false
}

// Cmp compares two UnsignedBigNumbers
func (n *UnsignedBigNumber) Cmp(other *UnsignedBigNumber) int {
	if n.Equals(other) {
		return 0
	}
	if n.LessThan(other) {
		return -1
	}
	return 1
}

// Sqrt calculates the square root of an UnsignedBigNumber
func (n *UnsignedBigNumber) Sqrt() *UnsignedBigNumber {
	if n.IsZero() {
		return &UnsignedBigNumber{Values: []byte{0}}
	}
	one := &UnsignedBigNumber{Values: []byte{1}}
	lo := one
	hi := n.Clone()
	for lo.LessThan(hi) {
		mid, _ := lo.Add(hi).Divide(one.Clone().Add(one))
		mid = mid.trim()
		if mid.Multiply(mid).LessThan(n) {
			lo = mid.Add(one)
		} else {
			hi = mid
		}
	}
	return lo.Subtract(one)
}

// ToBigInt converts the UnsignedBigNumber to a big.Int
func (n *UnsignedBigNumber) ToBigInt() *big.Int {
	bigInt := big.NewInt(0)
	bigInt.SetString(n.Hex(), 16)
	return bigInt
}

// Mod calculates the remainder of division of one UnsignedBigNumber by another
func (n *UnsignedBigNumber) Mod(other *UnsignedBigNumber) *UnsignedBigNumber {
	//start := time.Now()

	//_, remainder := n.Divide(other)
	bigInt := n.ToBigInt()
	otherBigInt := other.ToBigInt()

	bigInt.Mod(bigInt, otherBigInt)
	remainder := FromBigInt(bigInt)

	//fmt.Println("Mod:", time.Since(start))
	return remainder
}

// Multiply multiplies two UnsignedBigNumbers and returns the result
func (n *UnsignedBigNumber) Multiply(other *UnsignedBigNumber) *UnsignedBigNumber {
	//start := time.Now()
	var result *UnsignedBigNumber
	result = n.SimpleMultiply(other)
	//fmt.Println("Multiply:", time.Since(start))
	return result
}

// PowWithoutMod calculates the power of a number without modular arithmetic using binary exponentiation
func (n *UnsignedBigNumber) PowWithoutMod(exponent *UnsignedBigNumber) *UnsignedBigNumber {
	result := &UnsignedBigNumber{Values: []byte{1}}
	base := n.Clone()

	for !exponent.IsZero() {
		if exponent.Values[0]&1 == 1 {
			result = result.Multiply(base)
		}
		exponent = exponent.ShiftRight(1)
		base = base.Multiply(base)
	}

	return result
}

// Pow calculates the power of a number with modular arithmetic using binary exponentiation
func (n *UnsignedBigNumber) Pow(exponent, mod *UnsignedBigNumber) *UnsignedBigNumber {
	if mod == nil {
		return n.PowWithoutMod(exponent)
	}

	result := &UnsignedBigNumber{Values: []byte{1}}
	base := n.Clone().Mod(mod)
	exp := exponent.Clone()

	for !exp.IsZero() {
		if exp.Values[0]&1 == 1 {
			result = result.Multiply(base).Mod(mod)
		}
		exp = exp.ShiftRight(1)
		base = base.Multiply(base).Mod(mod)
	}

	return result
}

func (n *UnsignedBigNumber) GreaterThan(other *UnsignedBigNumber) bool {
	return !n.LessThan(other) && !n.Equals(other)
}
