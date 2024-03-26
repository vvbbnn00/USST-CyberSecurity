package bignumber

import "math/big"

// BigNumber represents a signed big number
type BigNumber struct {
	Sign   int // -1 for negative, 0 for zero, 1 for positive
	AbsVal *UnsignedBigNumber
}

// NewBigNumber creates a new BigNumber from a hex string
func NewBigNumber(hexStr string) *BigNumber {
	if len(hexStr) == 0 {
		return &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}
	}
	sign := 1
	if hexStr[0] == '-' {
		sign = -1
		hexStr = hexStr[1:]
	}
	return &BigNumber{Sign: sign, AbsVal: NewUnsignedBigNumber(hexStr)}
}

// Hex converts the BigNumber to a hex string
func (n *BigNumber) Hex() string {
	if n.Sign == 0 {
		return "0"
	}
	signStr := ""
	if n.Sign < 0 {
		signStr = "-"
	}
	return signStr + n.AbsVal.Hex()
}

// Add adds two BigNumbers and returns the result
func (n *BigNumber) Add(other *BigNumber) *BigNumber {
	if n.Sign == 0 {
		return other
	} else if other.Sign == 0 {
		return n
	} else if n.Sign == other.Sign {
		return &BigNumber{Sign: n.Sign, AbsVal: n.AbsVal.Add(other.AbsVal)}
	} else {
		if n.AbsVal.GreaterThan(other.AbsVal) {
			return &BigNumber{Sign: n.Sign, AbsVal: n.AbsVal.Subtract(other.AbsVal)}
		} else {
			return &BigNumber{Sign: -n.Sign, AbsVal: other.AbsVal.Subtract(n.AbsVal)}
		}
	}
}

// Subtract subtracts one BigNumber from another and returns the result
func (n *BigNumber) Subtract(other *BigNumber) *BigNumber {
	if n.Sign == 0 {
		return &BigNumber{Sign: -other.Sign, AbsVal: other.AbsVal}
	} else if other.Sign == 0 {
		return n
	} else if n.Sign == other.Sign {
		if n.AbsVal.GreaterThan(other.AbsVal) {
			return &BigNumber{Sign: n.Sign, AbsVal: n.AbsVal.Subtract(other.AbsVal)}
		} else {
			return &BigNumber{Sign: -n.Sign, AbsVal: other.AbsVal.Subtract(n.AbsVal)}
		}
	} else {
		return &BigNumber{Sign: n.Sign, AbsVal: n.AbsVal.Add(other.AbsVal)}
	}
}

// Multiply multiplies two BigNumbers and returns the result
func (n *BigNumber) Multiply(other *BigNumber) *BigNumber {
	if n.IsZero() || other.IsZero() {
		return &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}
	}
	sign := n.Sign * other.Sign
	return &BigNumber{Sign: sign, AbsVal: n.AbsVal.Multiply(other.AbsVal)}
}

// Divide divides one BigNumber by another and returns the result
func (n *BigNumber) Divide(other *BigNumber) (*BigNumber, *BigNumber) {
	if other.IsZero() {
		panic("division by zero")
	}
	if n.IsZero() {
		return &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}, &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}
	}
	sign := n.Sign * other.Sign
	quotient, remainder := n.AbsVal.Divide(other.AbsVal)
	return &BigNumber{Sign: sign, AbsVal: quotient}, &BigNumber{Sign: 1, AbsVal: remainder}
}

func (n *BigNumber) ToBigInt() *big.Int {
	bigInt := n.AbsVal.ToBigInt()
	if n.Sign == -1 {
		bigInt.Neg(bigInt)
	}
	return bigInt
}

// FromBigIntSigned converts a big.Int to a BigNumber
func FromBigIntSigned(bigInt *big.Int) *BigNumber {
	sign := 1
	if bigInt.Sign() == -1 {
		sign = -1
		bigInt.Neg(bigInt)
	}
	return &BigNumber{Sign: sign, AbsVal: FromBigInt(bigInt)}
}

// Mod calculates the modulus of one BigNumber by another and returns the result
func (n *BigNumber) Mod(other *BigNumber) *BigNumber {
	if other.IsZero() {
		panic("division by zero")
	}
	if n.IsZero() {
		return &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}
	}
	bigInt := n.ToBigInt()
	otherBigInt := other.ToBigInt()
	bigInt.Mod(bigInt, otherBigInt)
	return FromBigIntSigned(bigInt)
}

// Clone returns a copy of the BigNumber
func (n *BigNumber) Clone() *BigNumber {
	return &BigNumber{Sign: n.Sign, AbsVal: n.AbsVal.Clone()}
}

// GreaterThan returns true if n is greater than other, false otherwise
func (n *BigNumber) GreaterThan(other *BigNumber) bool {
	if n.Sign > other.Sign {
		return true
	} else if n.Sign < other.Sign {
		return false
	}
	if n.Sign == 0 {
		return false
	}
	return n.AbsVal.GreaterThan(other.AbsVal)
}

// LessThan returns true if n is less than other, false otherwise
func (n *BigNumber) LessThan(other *BigNumber) bool {
	if n.Sign < other.Sign {
		return true
	} else if n.Sign > other.Sign {
		return false
	}
	if n.Sign == 0 {
		return false
	}
	return n.AbsVal.LessThan(other.AbsVal)
}

// IsZero returns true if the BigNumber is zero, false otherwise
func (n *BigNumber) IsZero() bool {
	if n.AbsVal.IsZero() && n.Sign != 0 {
		n.Sign = 0
	}
	return n.Sign == 0
}

// Abs returns the absolute value of the BigNumber
func (n *BigNumber) Abs() *BigNumber {
	return &BigNumber{Sign: 1, AbsVal: n.AbsVal}
}

// Pow calculates n^exp mod m and returns the result
func (n *BigNumber) Pow(exp, m *BigNumber) *BigNumber {
	if m.Sign == 0 {
		panic("division by zero")
	}
	if exp.Sign == -1 {
		panic("negative exponent")
	}
	if exp.Sign == 0 {
		return &BigNumber{Sign: 1, AbsVal: NewUnsignedBigNumber("1")}
	}
	if n.Sign == 0 {
		return &BigNumber{Sign: 0, AbsVal: &UnsignedBigNumber{Values: []byte{0}}}
	}
	sign := n.Sign
	if sign == -1 {
		if exp.AbsVal.Values[0]&1 == 0 {
			sign = 1
		}
	}
	return &BigNumber{Sign: sign, AbsVal: n.AbsVal.Pow(exp.AbsVal, m.AbsVal)}
}

// Equals returns true if n is equal to other, false otherwise
func (n *BigNumber) Equals(other *BigNumber) bool {
	if n.Sign != other.Sign {
		return false
	}
	if n.Sign == 0 {
		return true
	}
	return n.AbsVal.Equals(other.AbsVal)
}
