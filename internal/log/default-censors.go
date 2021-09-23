package log

import "regexp"

var defaultCensorItems = []CensorItem{
	{
		SensitiveRegexp: regexp.MustCompile(`(?i)\n(Authorization: [a-z]+) ([^ ]+)\n`),
		Replacement:     "\n$1 <CENSORED>\n",
	},
}
