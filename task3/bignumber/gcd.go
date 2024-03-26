package bignumber

// GCDExt returns the greatest common divisor of n and m
func (n *BigNumber) GCDExt(m *BigNumber) (*BigNumber, *BigNumber, *BigNumber) {
	a, b := n.Abs(), m
	x0, y0 := NewBigNumber("1"), NewBigNumber("0")
	x1, y1 := NewBigNumber("0"), NewBigNumber("1")

	for !b.IsZero() {
		q, r := a.Divide(b)
		a, b = b, r
		x0, x1 = x1, x0.Subtract(q.Multiply(x1))
		y0, y1 = y1, y0.Subtract(q.Multiply(y1))
	}

	return a, x0, y0
}

// InvMod returns the modular inverse of n mod m
func (n *BigNumber) InvMod(m *BigNumber) *BigNumber {
	_, x, _ := n.GCDExt(m)
	return x.Mod(m)
}
