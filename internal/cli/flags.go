package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type flags struct {
	file string
	path string
}

const (
	defaultDataFilePath = "moods.json"
)

func ParseFlags() (flags, error) {

	file := flag.String("file", defaultDataFilePath, "moods saving file name in default location")
	path := flag.String("path", defaultDataFilePath, "full path to moods saving file")

	flag.Parse() //super important, otherwise flags won't work

	fileSet := false
	pathSet := false

	// flag.Visit iterates over flags that were provided by the user
	// the function passed to Visit is anonymous and is called once for each flag
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "file" {
			fileSet = true
		}

		if f.Name == "path" {
			pathSet = true
		}
	})

	if fileSet && pathSet {
		return flags{}, fmt.Errorf("flags file and path were used simultaneously, which is prohibited")
	}

	if pathSet {
		return flags{path: *path}, nil
	}

	return flags{file: *file}, nil

}

func BuildPath(parsedFlags flags) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	if parsedFlags.path != "" {
		return parsedFlags.path, nil
	}

	file := parsedFlags.file
	if file == "" {
		file = "moods.json"
	}

	return filepath.Join(configDir, "mindgo", file), nil
}
