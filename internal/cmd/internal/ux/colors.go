package ux

import "fmt"

type Color string

const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Cyan   Color = "\033[36m"
)

func Colored(text any, color Color) string {
	return fmt.Sprintf("%s%v%s", color, text, Reset)
}
