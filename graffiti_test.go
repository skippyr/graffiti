package graffiti

import (
	"testing"
)

func TestRemoveStyleAndCursorSequences(test *testing.T) {
	sample := "\x1b[31m[Here] Are Goats\x1b[0m"
	expectedResult := "[Here] Are Goats"
	result := removeStyleAndCursorSequences(sample)
	if result != expectedResult {
		test.Errorf(
			"Failed to remove style and cursor sequences. Expected \"%s\" (%d characters) but received \"%s\" (%d characters).",
			expectedResult,
			len(expectedResult),
			result,
			len(result),
		)
	}
}

func TestRemoveFormatSpecifiers(test *testing.T) {
	sample := "%F{red}Here %BAre%I %UGoats. %FMore %F}text %F{here.%r More"
	expectedResult := "Here Are Goats. More }text More"
	result := removeFormatSpecifiers(sample)
	if result != expectedResult {
		test.Errorf(
			"Failed to remove format specifiers. Expected \"%s\" (%d characters) but received \"%s\" (%d characters).",
			expectedResult,
			len(expectedResult),
			result,
			len(result),
		)
	}
}

