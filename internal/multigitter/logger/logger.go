package logger

import (
	"bufio"
	"io"
)

type logger interface {
	Infof(format string, args ...interface{})
}

// NewLogger creates a new logger that logs things that is written to the returned writer
func NewLogger(logger logger) io.WriteCloser {
	reader, writer := io.Pipe()

	// Print each line that is outputted by the script
	go func() {
		buf := bufio.NewReader(reader)
		for {
			line, err := buf.ReadString('\n')
			if line != "" {
				logger.Infof("Script output: %s", line)
			}
			if err != nil {
				return
			}
		}
	}()

	return writer
}
