package xmldiff

import (
	"fmt"
	"os"
)

const escape = "\x1b"

// Color chart: https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124

// https://no-color.org/
func isNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

func Bold(str string) string {
	if isNoColor() {
		return str
	}
	return fmt.Sprintf("%s[%sm%s%s[%sm", escape, "1;37", str, escape, "0")
}

func Red(str string) string {
	if isNoColor() {
		return str
	}
	return fmt.Sprintf("%s[%sm%s%s[%sm", escape, "0;31", str, escape, "0")
}

func Green(str string) string {
	if isNoColor() {
		return str
	}
	return fmt.Sprintf("%s[%sm%s%s[%sm", escape, "0;32", str, escape, "0")
}
