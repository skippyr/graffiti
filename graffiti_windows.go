package graffiti

import (
	"os"

	"golang.org/x/sys/windows"
)

func allowEscapingAnsiSequencesOnWindows() {
	// On Windows, most terminal emulators do not escape ANSI sequences by default.
	// Escaping is enabled by setting the ENABLE_VIRTUAL_TERMINAL_PROCESSING and ENABLE_PROCESSED_OUTPUT bits into the console mode value using an OR bitwise operation.
	// Here are some references that describe it:
	//     https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	//     https://learn.microsoft.com/en-us/windows/console/getconsolemode

	var consoleMode uint32
	stdoutHandle := windows.Handle(os.Stdout.Fd())
	if windows.GetConsoleMode(stdoutHandle, &consoleMode) != nil {
		return
	}
	isProcessingOutputAsAnsi := consoleMode&windows.ENABLE_PROCESSED_OUTPUT != 0
	if !isProcessingOutputAsAnsi {
		windows.SetConsolemode(stdoutHandle, consoleMode|windows.ENABLE_PROCESSED_OUTPUT)
	}
	isEscapingAnsiSequences := consoleMode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0
	if !isEscapingAnsiSequences {
		windows.SetConsoleMode(stdoutHandle, consoleMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}

func init() {
	allowEscapingAnsiSequencesOnWindows()
}
