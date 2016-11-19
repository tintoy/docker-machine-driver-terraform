package fetch

import (
	"github.com/hashicorp/go-getter"
)

// For now, we only support a subset of the sources that go-getter supports
// (weird problems with GitHub, for example).
var supportedGetters = map[string]getter.Getter{
	"file":  getter.Getters["file"],
	"http":  getter.Getters["http"],
	"https": getter.Getters["https"],
}

// For now, we only support a subset of the sources that go-getter supports
// (weird problems with GitHub, for example).
func supportedDetectors() []getter.Detector {
	return []getter.Detector{
		new(getter.FileDetector),
	}
}

// ParseSource parses the specified source into a form that Fetch can understand.
func ParseSource(source string, destination string) (string, error) {
	return getter.Detect(source, destination, supportedDetectors())
}

// Content downloads a URL into the given destination.
//
// destination must be a directory.
// If source is a file, it will be downloaded into destination
// with the basename of the URL.
// If source is a directory or archive, it will be unpacked directly into destination.
func Content(source string, destination string) error {
	return (&getter.Client{
		Src:     source,
		Dst:     destination,
		Mode:    getter.ClientModeAny,
		Getters: supportedGetters,
	}).Get()
}
