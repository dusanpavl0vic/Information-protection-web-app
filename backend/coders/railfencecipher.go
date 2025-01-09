package coders

func EncryptRailFence(message string, depth int) string {
	rails := make([][]rune, depth)
	for i := range rails {
		rails[i] = make([]rune, len(message))
		for j := range rails[i] {
			rails[i][j] = '\n'
		}
	}

	downward := false
	currentRow, currentCol := 0, 0

	for _, char := range message {
		if currentRow == 0 || currentRow == depth-1 {
			downward = !downward
		}
		rails[currentRow][currentCol] = char
		currentCol++
		if downward {
			currentRow++
		} else {
			currentRow--
		}
	}

	cipherText := make([]rune, 0, len(message))
	for i := 0; i < depth; i++ {
		for j := 0; j < len(message); j++ {
			if rails[i][j] != '\n' {
				cipherText = append(cipherText, rails[i][j])
			}
		}
	}

	return string(cipherText)
}

// Function to decrypt a message using Rail Fence Cipher
func decryptRailFence(cipher string, depth int) string {
	rails := make([][]rune, depth)
	for i := range rails {
		rails[i] = make([]rune, len(cipher))
		for j := range rails[i] {
			rails[i][j] = '\n'
		}
	}

	downward := false
	currentRow, currentCol := 0, 0

	for range cipher {
		if currentRow == 0 {
			downward = true
		}
		if currentRow == depth-1 {
			downward = false
		}
		rails[currentRow][currentCol] = '*'
		currentCol++
		if downward {
			currentRow++
		} else {
			currentRow--
		}
	}

	index := 0
	for i := 0; i < depth; i++ {
		for j := 0; j < len(cipher); j++ {
			if rails[i][j] == '*' && index < len(cipher) {
				rails[i][j] = rune(cipher[index])
				index++
			}
		}
	}

	plainText := make([]rune, 0, len(cipher))
	currentRow, currentCol = 0, 0
	for range cipher {
		if currentRow == 0 {
			downward = true
		}
		if currentRow == depth-1 {
			downward = false
		}
		if rails[currentRow][currentCol] != '*' {
			plainText = append(plainText, rails[currentRow][currentCol])
			currentCol++
		}
		if downward {
			currentRow++
		} else {
			currentRow--
		}
	}

	return string(plainText)
}
