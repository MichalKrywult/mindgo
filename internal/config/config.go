package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultDataFilePath = "moods.json"
)

type Config struct {
	DataFilePath string
}

type flags struct {
	file string
	path string
}

func parseFlags() (flags, error) {

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

func buildPath(parsedFlags flags) (string, error) {

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	var fullPath string

	if parsedFlags.path != "" {
		fullPath = parsedFlags.path
	} else {
		appDir := filepath.Join(configDir, "mindgo")
		fullPath = filepath.Join(appDir, parsedFlags.file)
	}

	return fullPath, nil
}

func GetConfig() (Config, error) {

	parsedFlags, err := parseFlags()
	if err != nil {
		return Config{}, err
	}

	fullPath, err := buildPath(parsedFlags)
	if err != nil {
		return Config{}, err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return Config{}, err
	}

	return Config{DataFilePath: fullPath}, nil
}
