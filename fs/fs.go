package fs

import (
	"regexp"
	"strings"
)

var (
	_invalidFileNameCharsRegex = regexp.MustCompile(`[\\/:*?"<>|]`)
)

func TrimInvalidFileNameChars(name string) string {
	return _invalidFileNameCharsRegex.ReplaceAllLiteralString(
		strings.Join(strings.Fields(name), " "),
		"_",
	)
}
