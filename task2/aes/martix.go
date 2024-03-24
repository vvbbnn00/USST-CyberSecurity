package aes

// mul2 GF(2^8) multiplication by 0x02
func mul2(x byte) byte {
	n := (x << 1) & 0xff
	if x&0x80 != 0 {
		n ^= 0x1b
	}
	return n
}

// mul3 GF(2^8) multiplication by 0x03
func mul3(x byte) byte {
	return x ^ mul2(x)
}

// mul9 GF(2^8) multiplication by 0x09
func mul9(x byte) byte {
	return x ^ mul2(mul2(mul2(x)))
}

// mul11 GF(2^8) multiplication by 0x0b
func mul11(x byte) byte {
	return x ^ mul2(x) ^ mul2(mul2(mul2(x)))
}

// mul13 GF(2^8) multiplication by 0x0d
func mul13(x byte) byte {
	return x ^ mul2(mul2(x)) ^ mul2(mul2(mul2(x)))
}

// mul14 GF(2^8) multiplication by 0x0e
func mul14(x byte) byte {
	return mul2(mul2(mul2(x))) ^ mul2(mul2(x)) ^ mul2(x)
}

// GFMul performs multiplication in GF(2^8)
func GFMul(x, mul byte) byte {
	switch mul {
	case 0x01:
		return x
	case 0x02:
		return mul2(x)
	case 0x03:
		return mul3(x)
	case 0x09:
		return mul9(x)
	case 0x0b:
		return mul11(x)
	case 0x0d:
		return mul13(x)
	case 0x0e:
		return mul14(x)
	default:
		panic("Invalid multiplier")
	}
}
