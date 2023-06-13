package graffiti

import (
	"testing"
	"strconv"
)

func TestRemoveAnsiEscapeSequences(test *testing.T) {
	for delimitersIterator, delimiter := range ansiEscapeSequencesDelimiters {
		sample := "\x1b[30" + string(delimiter) + strconv.Itoa(delimitersIterator)
		expectedResult := strconv.Itoa(delimitersIterator)
		result := removeAnsiEscapeSequences(&sample)
		if result != expectedResult {
			test.Errorf(
				"Failed to remove ANSI escape sequence with delimiter \"%c\". Expected \"%s\" (%d characters) but received \"%s\" (%d characters).",
				delimiter,
				expectedResult,
				len(expectedResult),
				result,
				len(result),
			)
		}
	}
}

