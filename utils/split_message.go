package utils

import "strings"

func SplitMessage(text string, maxLen int) []string {
	var parts []string

	for len(text) > 0 {
		if len(text) <= maxLen {
			parts = append(parts, text)
			break
		}

		// Ищем последний перенос строки перед maxLen
		splitAt := strings.LastIndex(text[:maxLen], "\n")
		if splitAt == -1 {
			splitAt = maxLen // Если переносов нет, режем по maxLen
		}

		parts = append(parts, text[:splitAt])
		text = text[splitAt:]
	}

	return parts
}
