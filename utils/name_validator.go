package utils

import (
	"regexp"
)

func IsValidFIO(input string) bool {
	// Регулярное выражение:
	// ^[А-ЯЁа-яё-]+       - Фамилия (обязательно)
	// (?:\s[А-ЯЁа-яё-]+)* - Имя и отчество (необязательно, можно 0, 1 или 2 слова)
	pattern := `^[А-ЯЁа-яё-]+(?:\s[А-ЯЁа-яё-]+)*$`
	matched, _ := regexp.MatchString(pattern, input)
	return matched
}
