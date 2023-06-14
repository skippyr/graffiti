package graffiti

import (
	"os"

	"golang.org/x/sys/windows"
)

func allowEscapingAnsiSequencesOnWindows() {
	// On Windows, most terminal emulators do not escape ANSI sequences by default.
	// Escaping is enabled by setting the ENABLE_VIRTUAL_TERMINAL_PROCESSING into the ConsoleMode bit.
	//     https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
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
