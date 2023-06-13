package graffiti

import (
	"fmt"
	"golang.org/x/term"
	"os"
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
	backgroundAnsiCode = 48
	boldAnsiCode       = 1
	foregroundAnsiCode = 38
	italicAnsiCode     = 4
	resetAnsiCode      = 0
	underlineAnsiCode  = 3
)
const (
	escapeCharacter                = '\x1b'
	formatSpecifierPrefixCharacter = '@'
	formatSpecifierOpenDelimiter   = '{'
	formatSpecifierCloseDelimiter  = '}'
)
const greatestFormatSpecifierValue = "magenta"

var ansiEscapeSequencesDelimiters = []rune{
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
var formatSpecifiers = map[rune][]int{
	'B': {boldAnsiCode, doNotExpectValue},
	'F': {foregroundAnsiCode, expectsValue},
	'I': {italicAnsiCode, doNotExpectValue},
	'K': {backgroundAnsiCode, expectsValue},
	'U': {underlineAnsiCode, doNotExpectValue},
	'r': {resetAnsiCode, doNotExpectValue},
}

func removeAnsiEscapeSequences(text *string) string {
	textWithoutStyleAndCursorSequences := ""
	hasUsedEscapeCharacter := false
	isEscaping := false
	for _, character := range *text {
		if isEscaping {
			for _, delimiter := range ansiEscapeSequencesDelimiters {
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
	valueLength := 0
	for _, character := range *text {
		if isReceivingValue {
			if character == ' ' || character == formatSpecifierCloseDelimiter || valueLength > len(greatestFormatSpecifierValue) {
				isReceivingValue = false
				valueLength = 0
			}
			valueLength++
			continue
		}
		if isExpectingValue == expectsValue {
			isExpectingValue = doNotExpectValue
			if character == formatSpecifierOpenDelimiter {
				isReceivingValue = true
				continue
			}
		}
		if character == formatSpecifierPrefixCharacter {
			isFormatting = !isFormatting
			if isFormatting {
				continue
			}
		}
		if isFormatting {
			isFormatting = false
			if len(formatSpecifiers[character]) > 0 {
				isExpectingValue = formatSpecifiers[character][1]
				continue
			}
		}
		textWithoutFormatSpecifiers = textWithoutFormatSpecifiers + string(character)
	}
	return textWithoutFormatSpecifiers
}

func treatText(stream int, text *string) string {
	treatedText := removeAnsiEscapeSequences(text)
	treatedText = removeFormatSpecifiers(&treatedText)
	if !term.IsTerminal(stream) {
		return removeFormatSpecifiers(&treatedText)
	}
	return treatedText
}

func writeToStream(stream int, text string, a ...any) (n int, err error) {
	treatedText := treatText(stream, &text)
	if stream == stdout {
		return fmt.Printf(treatedText, a...)
	}
	if stream == stderr {
		return fmt.Fprintf(os.Stderr, treatedText, a...)
	}
	return 0, nil
}

func Print(text string, a ...any) (n int, err error) {
	return writeToStream(stdout, text, a...)
}

func Println(text string, a ...any) (n int, err error) {
	return writeToStream(stdout, text+"\n", a...)
}

func EPrint(text string, a ...any) (n int, err error) {
	return writeToStream(stderr, text, a...)
}

func EPrintln(text string, a ...any) (n int, err error) {
	return writeToStream(stderr, text+"\n", a...)
}
