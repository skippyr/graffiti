package graffiti

import (
	"fmt"
	"os"
)

const (
	stdout = iota
	stderr
)

func treatText(text string) string {
	return text
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

func WriteLineToSdout(text string) {
	write(stdout, text+"\n")
}

func WriteToStderr(text string) {
	write(stderr, text)
}

func WriteLineToStderr(text string) {
	write(stderr, text+"\n")
}
