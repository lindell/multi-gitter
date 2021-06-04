package terminal

import "fmt"

// Link generates a link in that can be displayed in the terminal
// https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func Link(text, url string) string {
	return fmt.Sprintf("\x1B]8;;%s\a%s\x1B]8;;\a", url, text)
}

// Bold generates a bold text for the terminal
func Bold(text string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", text)
}
