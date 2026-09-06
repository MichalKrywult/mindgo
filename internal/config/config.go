package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DataFilePath string
}

const (
	defaultDataFilePath = "moods.json"
)

func GetConfig() (Config, error) {

	file := flag.String("file", defaultDataFilePath, "moods saving file name in default location")
	path := flag.String("path", defaultDataFilePath, "full path to moods saving file")

	flag.Parse() //super important, otherwise flags won't work

	if *file != defaultDataFilePath && *path != defaultDataFilePath {
		return Config{}, fmt.Errorf("flags file and path were used simultaneously, which is prohibited")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	var fullPath string

	if *path != defaultDataFilePath {
		fullPath = *path
	} else {
		appDir := filepath.Join(configDir, "mindgo")
		fullPath = filepath.Join(appDir, *file)
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return Config{}, err
	}

	return Config{DataFilePath: fullPath}, nil
}
