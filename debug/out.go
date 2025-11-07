package debug

import "fmt"

type DebugOutFunc func(string, ...any)

var Out DebugOutFunc = noop

func Verbose(enable bool) {
	if enable {
		Out = printf
	} else {
		Out = noop
	}
}

func noop(format string, args ...any) {}

func printf(format string, args ...any) {
	fmt.Printf(format, args...)
}
