package debug

type Context interface {
	// DebugPrintf writes a debug formatted message.
	// It doesn't do anything if debug mode is false.
	DebugPrintf(format string, a ...any)

	// DebugPrintf writes a debug message without newline.
	// It doesn't do anything if debug mode is false.
	DebugPrint(a ...any)

	// DebugPrintln writes a debug message with newline.
	// It doesn't do anything if debug mode is false.
	DebugPrintln(a ...any)
}
