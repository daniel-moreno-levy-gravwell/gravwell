package debug

import "fmt"

type DebugOutFunc func(string, ...any)

var DebugOut DebugOutFunc = Noop

func Noop(format string, args ...any) {}

func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}
