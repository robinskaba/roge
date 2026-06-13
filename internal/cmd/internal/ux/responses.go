package ux

import (
	"fmt"
	"os"
)

func Fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", Colored("fatal", Red), msg, err)
	os.Exit(1)
}

func Misuse(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", Colored("misuse", Yellow), msg)
	os.Exit(1)
}
