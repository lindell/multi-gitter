package log

import (
	"bytes"
	"strings"

	log "github.com/sirupsen/logrus"
)

// CensorFormatter makes sure sensitive data is not logged.
// It works as a middleware and sensors the data before sending it to an underlying formatter
type CensorFormatter struct {
	CensorItems         []CensorItem
	UnderlyingFormatter log.Formatter
}

// CensorItem is something that should be censored, Sensitive will be replaced with Replacement
type CensorItem struct {
	Sensitive   string
	Replacement string
}

// Format censors some data and sends the entry to the underlying formatter
func (f *CensorFormatter) Format(entry *log.Entry) ([]byte, error) {
	for _, s := range f.CensorItems {
		entry.Message = strings.ReplaceAll(entry.Message, s.Sensitive, s.Replacement)

		for key := range entry.Data {
			if str, ok := entry.Data[key].(string); ok {
				entry.Data[key] = strings.ReplaceAll(str, s.Sensitive, s.Replacement)
			}
			if bb, ok := entry.Data[key].([]byte); ok {
				entry.Data[key] = bytes.ReplaceAll(bb, []byte(s.Sensitive), []byte(s.Replacement))
			}
		}
	}
	return f.UnderlyingFormatter.Format(entry)
}
