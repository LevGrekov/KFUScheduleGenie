package utils

import (
	"strings"
)

// TranscribeName переводит русское имя и фамилию в английскую транслитерацию (например, "Агачев Юрий" → "Juriy", "Agachev")
func TranscribeName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName) // Разбиваем по пробелам
	if len(parts) < 2 {
		return "", "" // Если нет и имени, и фамилии
	}

	lastName = transliterate(parts[0])  // Фамилия (первое слово)
	firstName = transliterate(parts[1]) // Имя (второе слово)

	return firstName, lastName
}

// transliterate преобразует русскую строку в английскую транслитерацию
func transliterate(s string) string {
	translitMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
		'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
		'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "ju",
		'я': "ya",
	}

	var result strings.Builder
	for _, r := range strings.ToLower(s) {
		if en, ok := translitMap[r]; ok {
			result.WriteString(en)
		} else {
			result.WriteRune(r) // Оставляем как есть (например, цифры или латиницу)
		}
	}

	// Делаем первую букву заглавной (например, "agachev" → "Agachev")
	if result.Len() > 0 {
		return strings.Title(result.String())
	}
	return ""
}
