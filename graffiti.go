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
	backgroundAnsiCode   = 48
	boldAnsiCode         = 1
	foregroundAnsiCode   = 38
	italicAnsiCode       = 3
	resetAnsiCode        = 0
	underlineAnsiCode    = 4
	invalidAnsiColorCode = -1
	minimumAnsiColorCode = 0
	maximumAnsiColorCode = 1

	escapeCharacter                 = '\x1b'
	formatSpecifierPrefixCharacter  = '@'
	formatSpecifierOpenDelimiter    = '{'
	formatSpecifierCloseDelimiter   = '}'
	ansiEscapeSequenceOpenDelimiter = '['

	greatestFormatSpecifierValue = "magenta"
)

var ansiEscapeSequencesDelimiters = []rune{
	// A list of common ANSI escape sequences can be found at:
	//   https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	//   https://en.wikipedia.org/wiki/ANSI_escape_code
	'A', // Move cursor up
	'B', // Move cursor down
	'C', // Move cursor right
	'D', // Move cursor left
	'E', // Move cursor to beggining of next line
	'F', // Move cursor to beggining of previous line
	'H', // Move cursor
	'J', // Clear screen
	'K', // Clear line
	'M', // Move cursor one line up
	'm', // Style
	's', // Save cursor position
	'u', // Restore cursor to saved position
}
var threeBitsAnsiColors = []string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}
var formatSpecifiers = map[rune][]int{
	'B': {boldAnsiCode, doNotExpectValue},
	'F': {foregroundAnsiCode, expectsValue},
	'I': {italicAnsiCode, doNotExpectValue},
	'K': {backgroundAnsiCode, expectsValue},
	'U': {underlineAnsiCode, doNotExpectValue},
	'r': {resetAnsiCode, doNotExpectValue},
}

// Treats and returns a string with all ANSI escape sequences that can change styles, cursor position and clear contents removed.
func removeAnsiEscapeSequences(text *string) string {
	treatedText := ""
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
		isEscaping = hasUsedEscapeCharacter && character == ansiEscapeSequenceOpenDelimiter
		hasUsedEscapeCharacter = character == escapeCharacter
		if !isEscaping && !hasUsedEscapeCharacter {
			treatedText += string(character)
		}
	}
	return treatedText
}

func createStyleSequenceWithoutValue(ansiCode int) string {
	return fmt.Sprintf("%c[%dm", escapeCharacter, ansiCode)
}

func convertStringToAnsiColorCode(ansiColorAsString *string) int {
	for threeBitsAnsiColorCode, threeBitsAnsiColor := range threeBitsAnsiColors {
		if *ansiColorAsString == threeBitsAnsiColor {
			return threeBitsAnsiColorCode
		}
	}
	ansiColor, err := strconv.Atoi(*ansiColorAsString)
	if err != nil || ansiColor < minimumAnsiColorCode || ansiColor > maximumAnsiColorCode {
		return invalidAnsiColorCode
	}
	return ansiColor
}

func createStyleSequenceWithColor(ansiCode int, ansiColor int) string {
	return fmt.Sprintf("%c[%d;5;%dm", escapeCharacter, ansiCode, ansiColor)
}

func replaceFormatSpecifiers(text *string, isToAddNewLine bool) string {
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
				ansiColorCode := convertStringToAnsiColorCode(&value)
				if ansiColorCode != invalidAnsiColorCode {
					textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + createStyleSequenceWithColor(ansiCode, ansiColorCode)
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
	if isToAddNewLine {
		textWithFormatSpecifiersReplaced = textWithFormatSpecifiersReplaced + "\n"
	}
	return textWithFormatSpecifiersReplaced
}

func writeToStream(
	stream *os.File,
	text *string,
	isToAddNewLine bool,
	argumentsToFormat ...any,
) (bytesWritten int, err error) {
	*text = fmt.Sprintf(*text, argumentsToFormat...)
	*text = removeAnsiEscapeSequences(text)
	*text = replaceFormatSpecifiers(text, isToAddNewLine)
	if !term.IsTerminal(int(stream.Fd())) {
		*text = removeAnsiEscapeSequences(text)
	}
	return fmt.Fprint(stream, *text)
}

// Formats and prints a text to stdout. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Print(text string, argumentsToFormat ...any) (bytesWritten int, err error) {
	return writeToStream(os.Stdout, &text, false, argumentsToFormat...)
}

// Formats and prints a text to stdout with a new line character appended to its end. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Println(text string, argumentsToFormat ...any) (bytesWritten int, err error) {
	return writeToStream(os.Stdout, &text, true, argumentsToFormat...)
}

// Formats and prints a text to stderr. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Eprint(text string, argumentsToFormat ...any) (bytesWritten int, err error) {
	return writeToStream(os.Stderr, &text, false, argumentsToFormat...)
}

// Formats and prints a text to stderr with a new line character appended to its end. It accepts all format specifiers of fmt.Printf and also its own to deal with styling. It returns the number of bytes written and any write error encountered.
func Eprintln(text string, argumentsToFormat ...any) (bytesWritten int, err error) {
	return writeToStream(os.Stderr, &text, true, argumentsToFormat...)
}
