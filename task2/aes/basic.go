package aes

// subBytes substitutes the bytes in the matrix using the sbox
func subBytes(matrix Matrix, inv bool) Matrix {
	// determine the sbox
	box := sbox
	if inv {
		box = invSbox
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			value := matrix[i][j]
			// value: 0000 0000 1111 1111 -> row: 0000 0000, col: 1111 1111
			row := value >> 4
			col := value & 0x0f
			matrix[i][j] = box[row][col]
		}
	}
	return matrix
}

// shiftRows shifts the rows of the matrix
func shiftRows(matrix Matrix, inv bool) Matrix {
	result := new(Matrix)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			// x % 4 = x & 3
			if inv {
				// shift right
				result[i][j] = matrix[i][(j+4-i)&0x03]
			} else {
				// shift left
				result[i][j] = matrix[i][(j+i)&0x03]
			}
		}
	}

	return *result
}

// mixColumns mixes the columns of the matrix
func mixColumns(matrix Matrix, inv bool) Matrix {
	result := new(Matrix)
	mixColumnsMatrix := mixColumnsMatrix
	if inv {
		mixColumnsMatrix = invMixColumnsMatrix
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i][j] = GFMul(matrix[0][j], mixColumnsMatrix[i][0]) ^
				GFMul(matrix[1][j], mixColumnsMatrix[i][1]) ^
				GFMul(matrix[2][j], mixColumnsMatrix[i][2]) ^
				GFMul(matrix[3][j], mixColumnsMatrix[i][3])
		}
	}

	return *result
}

// addRoundKey adds the round key to the matrix
func addRoundKey(matrix Matrix, key Matrix) Matrix {
	result := new(Matrix)

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i][j] = matrix[i][j] ^ key[i][j]
		}
	}

	return *result
}

// EncryptBlock encrypts a block of data
func EncryptBlock(block Matrix, key *KeyExpansion) Matrix {
	// 128 -> 10, 192 -> 12, 256 -> 14
	totalRounds := int(key.keyType)*2 + 10

	// initial round
	block = addRoundKey(block, key.GetKey(0))

	// rounds
	for i := 1; i < totalRounds; i++ {
		//fmt.Println(block.Hex())

		block = subBytes(block, false)
		block = shiftRows(block, false)
		block = mixColumns(block, false)

		//fmt.Println(block.Hex())
		//getKey := key.GetKey(i)
		//fmt.Println(getKey.Hex())

		block = addRoundKey(block, key.GetKey(i))
	}

	// final round
	block = subBytes(block, false)
	block = shiftRows(block, false)
	block = addRoundKey(block, key.GetKey(totalRounds))

	return block
}

// DecryptBlock decrypts a block of data
func DecryptBlock(block Matrix, key *KeyExpansion) Matrix {
	totalRounds := int(key.keyType)*2 + 10

	// initial round
	block = addRoundKey(block, key.GetKey(totalRounds))
	block = shiftRows(block, true)
	block = subBytes(block, true)

	// rounds
	for i := totalRounds - 1; i > 0; i-- {
		block = addRoundKey(block, key.GetKey(i))
		block = mixColumns(block, true)
		block = shiftRows(block, true)
		block = subBytes(block, true)
	}

	// final round
	block = addRoundKey(block, key.GetKey(0))

	return block
}

// 474658a9 8c719558 0eb4d003 02ba2eb8
// 464658a8 8d719559 0fb4d002 03ba2eb9
