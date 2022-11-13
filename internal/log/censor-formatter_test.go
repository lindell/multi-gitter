package log

import (
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type dummyFormatter struct {
	entry logrus.Entry
}

func (f *dummyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.entry = *entry
	return []byte{}, nil
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name  string
		items []CensorItem
		entry logrus.Entry

		expectedMessage string
		expectedData    logrus.Fields
	}{
		{
			name: "simple",
			items: []CensorItem{
				{
					Sensitive:   "password",
					Replacement: "<censored>",
				},
			},
			entry: logrus.Entry{
				Message: "this is the password",
				Data: logrus.Fields{
					"string-field": "something password something",
					"bytes-field":  []byte("something password something"),
				},
			},
			expectedMessage: "this is the <censored>",
			expectedData: logrus.Fields{
				"string-field": "something <censored> something",
				"bytes-field":  []byte("something <censored> something"),
			},
		},
		{
			name: "multiple items",
			items: []CensorItem{
				{
					Sensitive:   "password",
					Replacement: "<censored-password>",
				},
				{
					Sensitive:   "token",
					Replacement: "<censored-token>",
				},
			},
			entry: logrus.Entry{
				Message: "a password and a token",
				Data: logrus.Fields{
					"string-field": "a password and a token",
					"bytes-field":  []byte("a password and a token"),
				},
			},
			expectedMessage: "a <censored-password> and a <censored-token>",
			expectedData: logrus.Fields{
				"string-field": "a <censored-password> and a <censored-token>",
				"bytes-field":  []byte("a <censored-password> and a <censored-token>"),
			},
		},
		{
			name: "http authorization",
			items: []CensorItem{
				{
					SensitiveRegexp: regexp.MustCompile(`(?i)\n(Authorization: [a-z]+) ([^ ]+)\n`),
					Replacement:     "\n$1 <censored>\n",
				},
			},
			entry: logrus.Entry{
				Data: logrus.Fields{
					"request": `GET /api/v4/ HTTP/1.1
Host: gitlab.com
User-Agent: Go-http-client/1.1
Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
Accept-Encoding: gzip

Some Data`,
					"request-as-byte-slice": []byte(`GET /api/v4/ HTTP/1.1
Host: gitlab.com
User-Agent: Go-http-client/1.1
Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==
Accept-Encoding: gzip

Some Data`),
				},
			},
			expectedData: logrus.Fields{
				"request": `GET /api/v4/ HTTP/1.1
Host: gitlab.com
User-Agent: Go-http-client/1.1
Authorization: Basic <censored>
Accept-Encoding: gzip

Some Data`,
				"request-as-byte-slice": []byte(`GET /api/v4/ HTTP/1.1
Host: gitlab.com
User-Agent: Go-http-client/1.1
Authorization: Basic <censored>
Accept-Encoding: gzip

Some Data`),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dummyFormatter := &dummyFormatter{}
			formatter := CensorFormatter{
				CensorItems:         test.items,
				UnderlyingFormatter: dummyFormatter,
			}
			_, err := formatter.Format(&test.entry)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedData, dummyFormatter.entry.Data)
			assert.Equal(t, test.expectedMessage, dummyFormatter.entry.Message)
		})
	}
}
