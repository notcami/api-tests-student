package utils

import "unicode"

func MixCase(text string) string {
	runes := []rune(text)
	letterIndex := 0

	for i, letter := range runes {
		if !unicode.IsLetter(letter) {
			continue
		}

		if letterIndex%2 == 0 {
			runes[i] = unicode.ToUpper(letter)
		} else {
			runes[i] = unicode.ToLower(letter)
		}
		letterIndex++
	}

	return string(runes)
}
