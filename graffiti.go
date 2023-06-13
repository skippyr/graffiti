package graffiti

import (
	"fmt"
	"strconv"
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
	backgroundAnsiCode = 48
	boldAnsiCode       = 1
	foregroundAnsiCode = 38
	italicAnsiCode     = 4
	resetAnsiCode      = 0
	underlineAnsiCode  = 3
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
	"black": 0,
	"red": 1,
	"green": 2,
	"yellow": 3,
	"blue": 4,
	"magenta": 5,
	"cyan": 6,
	"white": 7,
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

func createSimpleStyleSequence(ansiCode int) string {
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
					textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createSimpleStyleSequence(formatSpecifiers[character][0])
					hasStyle = true
				}
				continue
			}
		}
		textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + string(character)
	}
	if hasStyle {
		textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createSimpleStyleSequence(resetAnsiCode)
	}
	return textWithFormatSpecifiersReplaced
}

func treatText(stream int, text *string) string {
	treatedText := removeAnsiEscapeSequences(text)
	if !term.IsTerminal(stream) {
		return removeFormatSpecifiers(&treatedText)
	}
	treatedText = replaceFormatSpecifiers(&treatedText)
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
