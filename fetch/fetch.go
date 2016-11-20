package fetch

import (
	"os"
	"strings"

	"github.com/hashicorp/go-getter"
)

// For now, we only support a subset of the sources that go-getter supports
// (weird problems with GitHub, for example).
var supportedGetters = map[string]getter.Getter{
	"file":  getter.Getters["file"],
	"git":   getter.Getters["git"],
	"http":  getter.Getters["http"],
	"https": getter.Getters["https"],
}

// For now, we only support a subset of the sources that go-getter supports
// (weird problems with GitHub, for example).
func supportedDetectors() []getter.Detector {
	return []getter.Detector{
		new(getter.FileDetector),
		new(getter.GitHubDetector),
	}
}

// ParseSource parses the specified source into a form that Fetch can understand.
func ParseSource(source string) (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	isSourceDirectory := strings.HasSuffix(source, "/")
	parsedSource, err := getter.Detect(source, workingDirectory, supportedDetectors())
	if err != nil {
		return "", err
	}

	// Retain trailing slash if required.
	if isSourceDirectory && !strings.HasSuffix(parsedSource, "/") {
		parsedSource += "/"
	}

	return parsedSource, nil
}

// Content downloads a URL into the given destination.
//
// destination must be a directory.
// If source is a file, it will be downloaded into destination
// with the basename of the URL.
// If source is a directory or archive, it will be unpacked directly into destination.
func Content(source string, destination string) error {
	// Manual detection of source type (file / directory) since getter.ClientModeAny just assumes its a file.
	isDirectory := strings.HasSuffix(source, "/")
	clientMode := getter.ClientModeFile
	if isDirectory {
		clientMode = getter.ClientModeDir
		source = strings.TrimSuffix(source, "/") // Trim off the suffix since we only needed it for choosing the client mode
	}

	return (&getter.Client{
		Src:     source,
		Dst:     destination,
		Mode:    clientMode,
		Getters: supportedGetters,
	}).Get()
}
