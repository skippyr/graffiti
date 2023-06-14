package graffiti

import (
	"golang.org/x/sys/windows"
)

func init() {
	var consoleMode uint32
	stdoutHandle := windows.Handle(stdout)
	if windows.GetConsoleMode(stdoutHandle, &consoleMode) != nil {
		return
	}
	isEscapingAnsiSequences := consoleMode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0
	if !isEscapingAnsiSequences {
		windows.SetConsoleMode(stdoutHandle, consoleMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
