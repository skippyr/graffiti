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
const escapeCharacter = '\x1b'

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

var styleSequences = map[rune][]int {
	'B': {boldCode, doNotExpectValue},
	'F': {foregroundCode, expectsValue},
	'I': {italicCode, doNotExpectValue},
	'K': {backgroundCode, escapeCharacter},
	'U': {underlineCode, doNotExpectValue},
	'r': {resetCode, doNotExpectValue},
}

func removeStyleAndCursorSequences(text string) string {
	textWithoutStyleAndCursorSequences := ""
	hasUsedEscapeCharacter := false
	isEscaping := false
	for
		characters_iterator := 0;
		characters_iterator < len(text);
		characters_iterator++ {
		character := rune(text[characters_iterator])
		if isEscaping && character == 'm' {
			isEscaping = false
		}
		if hasUsedEscapeCharacter && character == '[' {
			isEscaping = true
		}
		if character == escapeCharacter {
			hasUsedEscapeCharacter = true
		} else {
			hasUsedEscapeCharacter = false
		}
		if !isEscaping {
			textWithoutStyleAndCursorSequences = textWithoutStyleAndCursorSequences + string(character)
		}
	}
	return textWithoutStyleAndCursorSequences
}

func getTextWithoutFormatSpecifiers(text string) string {
	return "(without) " + text
}

func treatText(stream int, text string) string {
	text = removeStyleAndCursorSequences(text)
	if !term.IsTerminal(stream) {
		return getTextWithoutFormatSpecifiers(text)
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
