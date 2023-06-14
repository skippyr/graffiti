package graffiti

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

const (
	expectsValue = iota
	doNotExpectValue
)
const (
	backgroundAnsiCode = 48
	boldAnsiCode       = 1
	foregroundAnsiCode = 38
	italicAnsiCode     = 3
	resetAnsiCode      = 0
	underlineAnsiCode  = 4
	invalidAnsiColor   = -1
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
var ansiColors = map[string]int{
	"black":   0,
	"red":     1,
	"green":   2,
	"yellow":  3,
	"blue":    4,
	"magenta": 5,
	"cyan":    6,
	"white":   7,
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

func createStyleSequenceWithoutValue(ansiCode int) string {
	return fmt.Sprintf("%c[%dm", escapeCharacter, ansiCode)
}

func convertStringToAnsiColor(value *string) int {
	if ansiColors[*value] != 0 || (ansiColors[*value] == 0 && *value == "black") {
		return ansiColors[*value]
	}
	ansiColor, err := strconv.Atoi(*value)
	if err != nil || ansiColor < 0 || ansiColor > 255 {
		return invalidAnsiColor
	}
	return ansiColor
}

func createStyleSequenceWithColor(ansiCode int, ansiColor int) string {
	return fmt.Sprintf("%c[%d;5;%dm", escapeCharacter, ansiCode, ansiColor)
}

func replaceFormatSpecifiers(text *string) string {
	textWithFormatSpecifiersReplaced := ""
	isFormatting := false
	isReceivingValue := false
	isExpectingValue := doNotExpectValue
	ansiCode := 0
	value := ""
	hasStyle := false
	for _, character := range *text {
		if isReceivingValue {
			if character == ' ' || character == formatSpecifierCloseDelimiter || len(value) > len(greatestFormatSpecifierValue) {
				ansiColor := convertStringToAnsiColor(&value)
				if ansiColor != invalidAnsiColor {
					textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createStyleSequenceWithColor(ansiCode, ansiColor)
					hasStyle = true
				}
				isReceivingValue = false
				value = ""
				ansiCode = 0
			} else {
				value = value + string(character)
			}
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
				ansiCode = formatSpecifiers[character][0]
				if isExpectingValue == doNotExpectValue {
					textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createStyleSequenceWithoutValue(formatSpecifiers[character][0])
					hasStyle = true
				}
				continue
			}
		}
		textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + string(character)
	}
	if hasStyle {
		textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createStyleSequenceWithoutValue(resetAnsiCode)
	}
	return textWithFormatSpecifiersReplaced
}

func writeToStream(stream *os.File, text string, a ...any) (n int, err error) {
	treatedText := fmt.Sprintf(text, a...)
	treatedText = removeAnsiEscapeSequences(&text)
	treatedText = replaceFormatSpecifiers(&treatedText)
	if !term.IsTerminal(int(stream.Fd())) {
		treatedText = removeAnsiEscapeSequences(&treatedText)
	}
	return fmt.Fprint(stream, treatedText)
}

// Formats and prints a text to stdout. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Print(text string, a ...any) (n int, err error) {
	return writeToStream(os.Stdout, text, a...)
}

// Formats and prints a text to stdout with a new line character appended to its end. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Println(text string, a ...any) (n int, err error) {
	return writeToStream(os.Stdout, text+"\n", a...)
}

// Formats and prints a text to stderr. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Eprint(text string, a ...any) (n int, err error) {
	return writeToStream(os.Stderr, text, a...)
}

// Formats and prints a text to stderr with a new line character appended to its end. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Eprintln(text string, a ...any) (n int, err error) {
	return writeToStream(os.Stderr, text+"\n", a...)
}
