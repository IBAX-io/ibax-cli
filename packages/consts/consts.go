package consts

import "strings"

// VERSION is current version
const VERSION = "1.0.0"

var BuildInfo string

func Version() string {
	return strings.TrimSpace(strings.Join([]string{VERSION, BuildInfo}, " "))
}
