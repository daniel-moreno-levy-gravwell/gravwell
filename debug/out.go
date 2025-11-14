package debug

import "fmt"

var enabled = false

func Out(format string, args ...any) {
	if enabled {
		fmt.Printf(format, args...)
	}
}

func SetPrintStdout(e bool) {
	enabled = e
}
