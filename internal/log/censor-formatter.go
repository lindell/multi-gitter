package log

import (
	"bytes"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NewCensorFormatter creates a new formatter that censors sensitive logs.
// It contains some default censoring rules, but additional items may be used
func NewCensorFormatter(underlyingFormatter log.Formatter, additionalCensoring ...CensorItem) *CensorFormatter {
	return &CensorFormatter{
		CensorItems:         append(defaultCensorItems, additionalCensoring...),
		UnderlyingFormatter: underlyingFormatter,
	}
}

// CensorFormatter makes sure sensitive data is not logged.
// It works as a middleware and sensors the data before sending it to an underlying formatter
type CensorFormatter struct {
	CensorItems         []CensorItem
	UnderlyingFormatter log.Formatter
}

// CensorItem is something that should be censored, Sensitive will be replaced with Replacement
type CensorItem struct {
	Sensitive       string
	SensitiveRegexp *regexp.Regexp
	Replacement     string
}

func (c CensorItem) stringReplace(str string) string {
	if c.Sensitive != "" {
		return strings.ReplaceAll(str, c.Sensitive, c.Replacement)
	}
	if c.SensitiveRegexp != nil {
		return c.SensitiveRegexp.ReplaceAllString(str, c.Replacement)
	}
	return str
}

func (c CensorItem) byteReplace(bb []byte) []byte {
	if c.Sensitive != "" {
		return bytes.ReplaceAll(bb, []byte(c.Sensitive), []byte(c.Replacement))
	}
	if c.SensitiveRegexp != nil {
		return c.SensitiveRegexp.ReplaceAll(bb, []byte(c.Replacement))
	}
	return bb
}

// Format censors some data and sends the entry to the underlying formatter
func (f *CensorFormatter) Format(entry *log.Entry) ([]byte, error) {
	for _, s := range f.CensorItems {
		entry.Message = s.stringReplace(entry.Message)

		for key := range entry.Data {
			if str, ok := entry.Data[key].(string); ok {
				entry.Data[key] = s.stringReplace(str)
			}
			if bb, ok := entry.Data[key].([]byte); ok {
				entry.Data[key] = s.byteReplace(bb)
			}
		}
	}
	return f.UnderlyingFormatter.Format(entry)
}
