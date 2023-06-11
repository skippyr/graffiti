package graffiti

import (
	"fmt"
	"os"
)

const (
	stdout = iota
	stderr
)
const escapeCharacter = '\x1b'

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

func treatText(text string) string {
	return removeStyleAndCursorSequences(text)
}

func write(stream int, text string) {
	treatedText := treatText(text)
	if stream == stdout {
		fmt.Print(treatedText)
	}
	if stream == stderr {
		fmt.Fprint(os.Stderr, treatedText)
	}
}

func WriteToStdout(text string) {
	write(stdout, text)
}

func WriteLineToStdout(text string) {
	write(stdout, text + "\n")
}

func WriteToStderr(text string) {
	write(stderr, text)
}

func WriteLineToStderr(text string) {
	write(stderr, text + "\n")
}
