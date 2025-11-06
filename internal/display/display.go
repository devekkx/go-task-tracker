package display

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	headerColor  = color.New(color.FgCyan, color.Bold)
	dimColor     = color.New(color.Faint)
	boldColor    = color.New(color.Bold)
)

// Success prints a success message
func Success(format string, args ...any) {
	successColor.Fprintf(os.Stdout, "✓ "+format+"\n", args...)
}

// Error prints an error message to stderr.
func Error(format string, args ...any) {
	errorColor.Fprintf(os.Stderr, "✗ "+format+"\n", args...)
}

// Warn prints a warning message.
func Warn(format string, args ...any) {
	warnColor.Fprintf(os.Stdout, "⚠ "+format+"\n", args...)
}

// Info prints an informational message.
func Info(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}
