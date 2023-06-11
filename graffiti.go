package graffiti

import (
	"fmt"
)

func createStyleSequence(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

func style(text string) string {
	simpleFormatSpecifiers := map[rune]int {
		'B': 1, // Bold
		'I': 3, // Italic
		'U': 4, // Underline
		'r': 0, // Reset
	}
	isFormatting := false
	hasStyle := false
	styledText := ""
	for
		characters_iterator := 0;
		characters_iterator < len(text);
		characters_iterator ++ {
		character := rune(text[characters_iterator])
		if character == '%' {
			if isFormatting {
				isFormatting = false
			} else {
				isFormatting = true
			}
		}
		if !isFormatting {
			styledText = styledText + string(character)
			continue
		}
		for formatSpecifier, formatCode := range simpleFormatSpecifiers {
			if character == formatSpecifier {
				styledText = styledText + createStyleSequence(formatCode)
				isFormatting = false
				hasStyle = true
			}
		}
	}
	if hasStyle {
		styledText = styledText + createStyleSequence(0)
	}
	return styledText
}

func StdoutWrite(text string) {
	styledText := style(text)
	fmt.Print(styledText)
}

