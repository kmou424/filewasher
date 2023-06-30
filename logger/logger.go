package logger

import "fmt"

const tag = "FileWasher"

func Logf(format string, msg ...string) {
	if len(msg) == 0 {
		fmt.Printf(fmt.Sprintf("%s: %s\n", tag, format))
		return
	}
	fmt.Printf(fmt.Sprintf("%s: %s\n", tag, format), msg)
}
