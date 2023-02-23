package helper

func TruncateString(s string, maxLength int) string {
	maxLength = maxLength - 3
	runes := []rune(s)
	if len(runes) > maxLength {
		return string(runes[:maxLength]) + "..."
	}
	return s
}
