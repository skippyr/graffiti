package graffiti

import (
	"testing"
)

func TestRemoveStyleAndCursorSequences(t *testing.T) {
	sample := "\x1b[31mHere Are Goats\x1b[0m"
	expectedResult := "Here Are Goats"
	result := removeStyleAndCursorSequences(sample)
	if result != expectedResult {
		t.Errorf(
			"Failed to remove style and cursor sequences. Expect \"%s\" (%d characters) but received \"%s\" (%d characters).",
			expectedResult,
			len(expectedResult),
			result,
			len(result),
		)
	}
}
