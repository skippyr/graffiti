package graffiti

import (
	"os"

	"golang.org/x/sys/windows"
)

func allowEscapingAnsiSequencesOnWindows() {
	// On Windows, most terminal emulators do not escape ANSI sequences by default.
	// Escaping is enabled by setting the ENABLE_VIRTUAL_TERMINAL_PROCESSING bit into the console mode value using an OR bitwise operation.
	//     https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	//     https://learn.microsoft.com/en-us/windows/console/getconsolemode
	//     https://pkg.go.dev/golang.org/x/sys/windows
	var consoleMode uint32
	stdoutHandle := windows.Handle(os.Stdout.Fd())
	if windows.GetConsoleMode(stdoutHandle, &consoleMode) != nil {
		return
	}
	isEscapingAnsiSequences := consoleMode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0
	if !isEscapingAnsiSequences {
		windows.SetConsoleMode(stdoutHandle, consoleMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}

func init() {
	allowEscapingAnsiSequencesOnWindows()
}
