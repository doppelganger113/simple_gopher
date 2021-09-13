package simple_gopher

import "unicode"

var connectionsRunes = [3]rune{'-', '_', ' '}

func isConnectionChar(char rune) bool {
	for _, value := range connectionsRunes {
		if value == char {
			return true
		}
	}
	return false
}

func areRestConnectionChars(index int, text string) bool {
	if index >= len(text) {
		return false
	}
	restOfTheText := text[index:]
	for _, char := range restOfTheText {
		if isConnectionChar(char) == false {
			return false
		}
	}

	return true
}

// FormatForSeo Creates an SEO friendly text from provided, example:
// from: _some1 ran----dom __test -- to add-5-
// to: some1-ran-dom-test-to-add-5
func FormatForSeo(name string) string {
	var newText []rune
	var lastChar rune

	for i, char := range name {
		// skip unknown characters
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !isConnectionChar(char) {
			continue
		}

		isFirstOrLast := i == 0 || i == len(name)-1

		if isFirstOrLast && (char == ' ' || char == '_' || char == '-') {
			lastChar = char
			continue
		} else if char == '_' {
			if areRestConnectionChars(i, name) {
				return string(newText)
			}
			if isConnectionChar(lastChar) {
				lastChar = '_'
				continue
			}
			newText = append(newText, '-')
		} else if char == '-' {
			if areRestConnectionChars(i, name) {
				return string(newText)
			}
			if isConnectionChar(lastChar) {
				lastChar = '-'
				continue
			}
			newText = append(newText, '-')
		} else if char == ' ' {
			if areRestConnectionChars(i, name) {
				return string(newText)
			}
			if isConnectionChar(lastChar) {
				lastChar = char
				continue
			}
			newText = append(newText, '-')
		} else {
			newText = append(newText, char)
		}

		lastChar = char
	}

	return string(newText)
}
