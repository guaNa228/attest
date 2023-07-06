package translit

import (
	"strings"
	"unicode"
)

var RussianLogin = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "yo",
	'ж': "zh",
	'з': "z",
	'и': "i",
	'й': "j",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "h",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "sch",
	'ъ': "",
	'ы': "y",
	'ь': "",
	'э': "e",
	'ю': "ju",
	'я': "ya",
}

var RussianEmail = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "_",
	'ж': "zgh",
	'з': "z",
	'и': "i",
	'й': "y",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "h",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "sh",
	'ъ': "",
	'ы': "u",
	'ь': "",
	'э': "e",
	'ю': "u",
	'я': "ya",
}

func ToLatin(s string, RussianASCII map[rune]string) string {
	runes := []rune(s)
	out := make([]rune, 0, len(s))
	for i, currentRune := range runes {
		if tr, ok := RussianASCII[unicode.ToLower(currentRune)]; ok {
			if tr == "" {
				continue
			}
			if unicode.IsUpper(currentRune) {
				// Correctly translate case of successive characters:
				// ЩИ -> SCHI
				// Щи -> Schi
				if i+1 < len(runes) && !unicode.IsUpper(runes[i+1]) {
					t := []rune(tr)
					t[0] = unicode.ToUpper(t[0])
					out = append(out, t...)
					continue
				}
				out = append(out, []rune(strings.ToUpper(tr))...)
				continue
			}
			out = append(out, []rune(tr)...)
		} else {
			out = append(out, currentRune)
		}
	}
	return string(out)
}
