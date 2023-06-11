package graffiti

import (
	"fmt"
)

func create_style_sequence(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

func style(text string) string {
	simple_format_specifiers := map[rune]int {
		'B': 1, // Bold
		'I': 3, // Italic
		'U': 4, // Underline
	}
	is_formatting := false
	has_style := false
	styled_text := ""
	for
		characters_iterator := 0;
		characters_iterator < len(text);
		characters_iterator ++ {
		character := rune(text[characters_iterator])
		if character == '%' {
			if is_formatting {
				is_formatting = false
			} else {
				is_formatting = true
			}
		}
		if !is_formatting {
			styled_text = styled_text + string(character)
		}
		for format_specifier, format_code := range simple_format_specifiers {
			if character == format_specifier {
				styled_text = styled_text + create_style_sequence(format_code)
				is_formatting = false
				has_style = true
			}
		}
	}
	if has_style {
		styled_text = styled_text + create_style_sequence(0)
	}
	return styled_text
}

func StdoutWrite(text string) {
	styledText := style(text)
	fmt.Print(styledText)
}

