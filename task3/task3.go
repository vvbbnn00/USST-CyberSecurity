package main

import "USST-CyberSecurity/task3/bignumber"

func main() {
	// 1024-bit = 128 bytes
	//p := bignumber.RandPrime(128)
	//q := bignumber.RandPrime(128)
	//println(p.Hex())
	//println(q.Hex())

	p := bignumber.NewBigNumber("f26b62487df0f6d6042dff1d155391281e906e78a2fbd48a720c75e689be4335c0b50bf092fe09a86b7ddea94203e440e01116ed96f836f6374bdcadf5130265c771e0f1cebc6cffcf78b522335587e688509f7e7d2d55dbf3341d495c21f6ca26632e3bc33f6bc28ba2dcaa8644a373d72fb7e6d2f69af5d363fe9527308adb")
	q := bignumber.NewBigNumber("d89b3d9fae9850bea9cbc7ea66b73b7bb0e2b12411e87d52707df8a4f50087c352d9cf576d609f6e848bf382e3ec32efc0d4f407620dd87bec4921e8d2b544419b476a3a9e15dfd50125c648f1e2e2f370339f8ba56d100bc7f3430b751ecbcf1b0d099b0f237f2b02dd13a1b6632a8e1904e6c59e6735de81712e76cde8e657")
	n := p.Multiply(q)

	println(n.Hex())

	// e is a small odd integer
	e := bignumber.NewBigNumber("10001")

	// phi(n) = (p-1)(q-1)
	p1 := p.Subtract(bignumber.NewBigNumber("1"))
	q1 := q.Subtract(bignumber.NewBigNumber("1"))

	phi := p1.Multiply(q1)
	println(phi.Hex())

	// d = e^-1 mod phi
	d := e.InvMod(phi)
	println(d.Hex())

	message := bignumber.NewBigNumber("1234567890")
	ciphertext := message.Pow(e, n)
	println(ciphertext.Hex())

	decrypted := ciphertext.Pow(d, n)
	println(decrypted.Hex())
}
