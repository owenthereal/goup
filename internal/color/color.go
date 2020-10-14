package color

import "fmt"

func Str2Red(str string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", str)
}

func Str2Cyan(str string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", str)
}
