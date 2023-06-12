package graffiti

import (
	"testing"
	"strconv"
)

func TestRemoveHiddenSequences(test *testing.T) {
	sample := ""
	expectedResult := ""
	for
		delimitersIterator := 0;
		delimitersIterator < len(hiddenSequencesDelimiters);
		delimitersIterator ++ {
		delimiter := hiddenSequencesDelimiters[delimitersIterator]
		sample = sample + " \x1b[30" + string(delimiter) + strconv.Itoa(delimitersIterator)
		expectedResult = expectedResult + " " + strconv.Itoa(delimitersIterator)
	}
	result := removeHiddenSequences(sample)
	if result != expectedResult {
		test.Errorf(
			"Failed to remove hidden sequences. Expected \"%s\" (%d characters) but received \"%s\" (%d characters).",
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

