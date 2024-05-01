package terminal

import "fmt"

// Printer formats things to the terminal
type Printer struct {
	Plain bool // Don't use terminal formatting
}

var DefaultPrinter = &Printer{}

// Link generates a link in that can be displayed in the terminal
// https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func (t *Printer) Link(text, url string) string {
	if t.Plain {
		return text
	}

	return fmt.Sprintf("\x1B]8;;%s\a%s\x1B]8;;\a", url, text)
}

// Bold generates a bold text for the terminal
func (t *Printer) Bold(text string) string {
	if t.Plain {
		return text
	}

	return fmt.Sprintf("\033[1m%s\033[0m", text)
}

// Link generates a link in that can be displayed in the terminal using the default terminal printer
func Link(text, url string) string {
	return DefaultPrinter.Link(text, url)
}

// Bold generates a bold text for the terminal using the default terminal printer
func Bold(text string) string {
	return DefaultPrinter.Bold(text)
}
