package graffiti

import (
	"testing"
	"strconv"
)

func TestRemoveAnsiEscapeSequences(test *testing.T) {
	sample := ""
	expectedResult := ""
	for delimitersIterator, delimiter := range ansiEscapeSequencesDelimiters {
		sample = sample + " \x1b[30" + string(delimiter) + strconv.Itoa(delimitersIterator)
		expectedResult = expectedResult + " " + strconv.Itoa(delimitersIterator)
	}
	result := removeAnsiEscapeSequences(&sample)
	if result != expectedResult {
		test.Errorf(
			"Failed to remove ANSI escape sequences. Expected \"%s\" (%d characters) but received \"%s\" (%d characters).",
			expectedResult,
			len(expectedResult),
			result,
			len(result),
		)
	}
}

