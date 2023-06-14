package graffiti

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

type FormatSpecifier struct {
	character        rune
	ansiStyleCode    int
	doesExpectsColor bool
}

const (
	resetAnsiStyleCode   = 0
	invalidAnsiCode      = -1
	minimumAnsiColorCode = 0
	maximumAnsiColorCode = 255

	escapeCharacter                 = '\x1b'
	formatSpecifierPrefixCharacter  = '@'
	formatSpecifierOpenDelimiter    = '{'
	formatSpecifierCloseDelimiter   = '}'
	ansiEscapeSequenceOpenDelimiter = '['

	greatestFormatSpecifierColor = "magenta"
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
var formatSpecifiers = []FormatSpecifier{
	{
		// Bold
		character:        'B',
		ansiStyleCode:    1,
		doesExpectsColor: false,
	},
	{
		// Italic
		character:        'I',
		ansiStyleCode:    3,
		doesExpectsColor: false,
	},
	{
		// Underline
		character:        'U',
		ansiStyleCode:    4,
		doesExpectsColor: false,
	},
	{
		// Foreground
		character:        'F',
		ansiStyleCode:    38,
		doesExpectsColor: true,
	},
	{
		// Background
		character:        'K',
		ansiStyleCode:    48,
		doesExpectsColor: true,
	},
	{
		// Reset
		character:        'r',
		ansiStyleCode:    resetAnsiStyleCode,
		doesExpectsColor: false,
	},
}

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

func convertStringToAnsiColorCode(ansiColorAsString *string) int {
	for threeBitsAnsiColorCode, threeBitsAnsiColor := range threeBitsAnsiColors {
		if *ansiColorAsString == threeBitsAnsiColor {
			return threeBitsAnsiColorCode
		}
	}
	ansiColor, err := strconv.Atoi(*ansiColorAsString)
	if err != nil || ansiColor < minimumAnsiColorCode || ansiColor > maximumAnsiColorCode {
		return invalidAnsiCode
	}
	return ansiColor
}

func createAnsiStyleSequence(ansiStyleCode int) string {
	return fmt.Sprintf("%c[%dm", escapeCharacter, ansiStyleCode)
}

func createAnsiStyleSequenceWithColor(ansiStyleCode int, ansiColor int) string {
	return fmt.Sprintf("%c[%d;5;%dm", escapeCharacter, ansiStyleCode, ansiColor)
}

func replaceFormatSpecifiers(text *string, isToAddNewLine bool) string {
	ansiStyleCode := invalidAnsiCode
	color := ""
	hasStyle := false
	isExpectingColor := false
	isFormatting := false
	isReadingColor := false
	treatedText := ""
	for _, character := range *text {
		if isReadingColor {
			if character == ' ' || character == formatSpecifierCloseDelimiter || len(color) > len(greatestFormatSpecifierColor) {
				ansiColorCode := convertStringToAnsiColorCode(&color)
				if ansiColorCode != invalidAnsiCode {
					treatedText += createAnsiStyleSequenceWithColor(ansiStyleCode, ansiColorCode)
					hasStyle = true
				}
				isReadingColor = false
				color = ""
				ansiStyleCode = invalidAnsiCode
			} else {
				color += string(character)
			}
			continue
		}
		if isExpectingColor {
			isExpectingColor = false
			if character == formatSpecifierOpenDelimiter {
				isReadingColor = true
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
			isFormatSpecifier := false
			for _, formatSpecifier := range formatSpecifiers {
				if character == formatSpecifier.character {
					isFormatSpecifier = true
					isExpectingColor = formatSpecifier.doesExpectsColor
					ansiStyleCode = formatSpecifier.ansiStyleCode
					if !isExpectingColor {
						treatedText += createAnsiStyleSequence(formatSpecifier.ansiStyleCode)
						hasStyle = true
					}
					break
				}
			}
			if isFormatSpecifier {
				continue
			}
		}
		treatedText += string(character)
	}
	if hasStyle {
		treatedText += createAnsiStyleSequence(resetAnsiStyleCode)
	}
	if isToAddNewLine {
		treatedText += "\n"
	}
	return treatedText
}

func writeToStream(stream *os.File, text *string, isToAddNewLine bool, argumentsToFormat ...any) (bytesWritten int, err error) {
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
