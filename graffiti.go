package graffiti

import (
	"fmt"
	"os"
	"golang.org/x/term"
)

const (
	stdout = iota + 1
	stderr
)
const (
	expectsValue = iota
	doNotExpectValue
)
const (
	backgroundCode = 48
	boldCode = 1
	foregroundCode = 38
	italicCode = 4
	resetCode = 0
	underlineCode = 3
)
const (
	escapeCharacter = '\x1b'
	formatSpecifierPrefixCharacter = '@'
	formatSpecifierOpenDelimiter = '{'
	formatSpecifierCloseDelimiter = '}'
)
const greatestFormatSpecifierValue = len("magenta")
var hiddenSequencesDelimiters = []rune {
	'H', // Move cursor
	'J', // Clear screen
	'A', // Move cursor up
	'B', // Move cursor down
	'C', // Move cursor right
	'D', // Move cursor left
	'E', // Move cursor to beggining of next line
	'F', // Move cursor to beggining of previous line
	'm', // Style
}
var formatSpecifiers = map[rune][]int {
	'B': {boldCode, doNotExpectValue},
	'F': {foregroundCode, expectsValue},
	'I': {italicCode, doNotExpectValue},
	'K': {backgroundCode, expectsValue},
	'U': {underlineCode, doNotExpectValue},
	'r': {resetCode, doNotExpectValue},
}

func removeHiddenSequences(text *string) string {
	textWithoutStyleAndCursorSequences := ""
	hasUsedEscapeCharacter := false
	isEscaping := false
	for _, character := range *text {
		if isEscaping {
			for
				delimitersIterator := 0;
				delimitersIterator < len(hiddenSequencesDelimiters);
				delimitersIterator ++ {
				delimiter := hiddenSequencesDelimiters[delimitersIterator]
				if character == delimiter {
					isEscaping = false
					break
				}
			}
			continue
		}
		isEscaping = hasUsedEscapeCharacter && character == '['
		hasUsedEscapeCharacter = character == escapeCharacter
		if hasUsedEscapeCharacter || isEscaping {
			continue
		}
		textWithoutStyleAndCursorSequences = textWithoutStyleAndCursorSequences + string(character)
	}
	return textWithoutStyleAndCursorSequences
}

func removeFormatSpecifiers(text *string) string {
	textWithoutFormatSpecifiers := ""
	isFormatting := false
	isExpectingValue := doNotExpectValue
	isReceivingValue := false
	valueSize := 0
	for _, character := range *text {
		if isReceivingValue {
			valueSize ++
			if
				character == ' ' ||
				character == formatSpecifierCloseDelimiter ||
				valueSize > greatestFormatSpecifierValue {
				isReceivingValue = false
				valueSize = 0
			}
			continue
		}
		if character == formatSpecifierPrefixCharacter {
			isFormatting = !isFormatting
			continue
		}
		if isFormatting {
			if len(formatSpecifiers[character]) != 0 {
				isExpectingValue = formatSpecifiers[character][1]
			}
			isFormatting = false
			continue
		}
		if isExpectingValue == expectsValue {
			isExpectingValue = doNotExpectValue
			if character == formatSpecifierOpenDelimiter {
				isReceivingValue = true
				continue
			}
		}
		textWithoutFormatSpecifiers = textWithoutFormatSpecifiers + string(character)
	}
	return textWithoutFormatSpecifiers
}

func treatText(stream int, text string) string {
	text = removeHiddenSequences(&text)
	if !term.IsTerminal(stream) {
		return removeFormatSpecifiers(&text)
	}
	return text
}

func writeToStream(stream int, text string) {
	treatedText := treatText(stream, text)
	if stream == stdout {
		fmt.Print(treatedText)
	}
	if stream == stderr {
		fmt.Fprint(os.Stderr, treatedText)
	}
}

func Write(text string) {
	writeToStream(stdout, text)
}

func WriteLine(text string) {
	writeToStream(stdout, text + "\n")
}

func ErrWrite(text string) {
	writeToStream(stderr, text)
}

func ErrWriteLine(text string) {
	writeToStream(stderr, text + "\n")
}
