package util

// FilterWords filters out non-alphabetic characters from the text
func FilterWords(text string) string {
	filtered := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			filtered += string(char)
		}
	}
	return filtered
}
